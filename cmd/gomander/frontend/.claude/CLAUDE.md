# CLAUDE.md

The Gomander frontend: React 19, TypeScript, Vite, Zustand, TailwindCSS v4,
shadcn/ui, react-hook-form + Zod, Biome, pnpm. It talks to the Go backend through
Wails bindings. See [the root CLAUDE.md](../../../../.claude/CLAUDE.md) for the backend.

## Where the reasoning lives

- **[CONTEXT.md](../../../../CONTEXT.md)** — the glossary, shared with the backend.
- **[docs/adr/](../../../../docs/adr/)** — decisions and the alternatives they rejected.
  ADR-0001 (ports at the third-party edges), 0002 (autosave), 0003 (browser storage)
  and 0005 (coverage) all bind code in this directory.

`pnpm run` lists the scripts. The dev server runs from the repo root (`make dev`).

## Layout

```
src/
├── contracts/         # the backend seam - the ONLY place wailsjs is imported
│   ├── ports.ts       # the contract, hand-declared
│   ├── service.ts     # the services the app calls, plus the test-only swap
│   └── adapters/      # wails.ts (production), inMemory.ts (tests)
├── commandOutput/     # a command's output, from Wails event to pixels
│   ├── commandOutput.ts
│   ├── ports.ts       # the terminal emulator contract
│   └── adapters/      # xterm.ts (production), recording.ts (tests)
├── useCases/          # write operations, by domain
├── queries/           # read-only fetches
├── store/             # one Zustand store per domain
├── screens/ components/ design-system/
├── hooks/ contexts/ helpers/ types/ constants/
└── testing/           # fake backend and terminal installers, builders
```

`@/` maps to `src/`. Routing is `HashRouter`, not `BrowserRouter` — a Wails requirement.

## Rules

### Backend calls

Adding one is four steps: declare it on the right type in `contracts/ports.ts`, map it
in `adapters/wails.ts`, implement it in `adapters/inMemory.ts`, then call it from a use
case. Biome fails any `wailsjs` import outside `contracts/`, and any `@xterm/*` import
outside `commandOutput/`.

### Mutations own their outcome

A mutation in `useCases/` owns the whole thing: the backend call, the stores it has to
refresh, its own `i18n.t` toast on success and on failure, and any state the change
forces (disposing a terminal, clearing the active command). It never throws — it returns
`false`, so a caller with UI to run on success can branch on it.

The call site keeps only what is local to the UI: closing a modal, resetting a form,
`stopPropagation`, selecting the clicked command. A `try`/`catch`, a `fetchX()` or a
`parseError` around a mutation means that mutation has not absorbed its outcome yet.

Refresh through `refreshAfterMutation(...queries)` — the outcome is already reported by
then, so a failing refresh must not reject. Equivalent mutations refresh the same set: a
command's create and delete refresh the commands and the groups that name it, while an
edit refreshes the commands alone because a group holds no copy to update; a group's own
mutations refresh the groups alone; a project mutation that changes which projects exist
refreshes the available ones, and an edit refreshes the opened project too.
`closeProject` and `exportProject` refresh nothing.

`createProject` is the one exception left — still a `dataService` call from
`CreateProjectModal` with an `onSuccess` prop.

### Modals

`components/modals/Command/common/formMapping.ts` maps between a form and a `Command`,
including the one-error-pattern-per-line encoding. A modal that splits or joins that
string itself should be calling this instead.

Every create/edit modal reads the same way: the mutation reports its own outcome and
returns a boolean, the modal returns early on `false`, and on success closes and resets
in that order — `setOpen(false)`, then `form.reset()`. A rejected save leaves the modal
open with what the user typed still in it.

### Modules that own their sequencing

Call these; do not rebuild what they do.

- `commandOutput/commandOutput.ts` — buffering, the 30 ms flush, timestamping, the
  backfill for terminals not yet on screen, the bounded tail and the reset ordering.
- `hooks/useAutosavedForm.ts` — the dirty check, the debounce, the submit, the pending
  flag. A call site running its own `setTimeout` or diffing against a `useRef` should be
  using this. The rules live in `createAutosaver`, a plain factory unit-tested with no
  renderer.
- `store/sidebarSections.ts` — `useSidebarSection(id)` hands back a `useState`-shaped
  pair; storage keys and encoding stay inside.

### State and events

Components read state with `useCommandStore(state => state.field)`; use cases reach it
with `commandStore.getState()` / `.setState()`. Keep selectors narrow.

Every backend event listener is registered in `EventListenersContainer`, which is wiring
only: one event, one call, no pipeline of its own. New events are declared in
`contracts/types.ts`.

## Testing

Vitest, jsdom, tests next to the code as `*.test.ts(x)`. Node ≥22.10 — jsdom 30 pulls an
undici that calls `worker_threads.markAsUncloneable`, and on older Node every test file
fails to start.

- `installInMemoryBackend()` (from `@/testing/backend.ts`) swaps every service;
  `resetBackendServices()` puts Wails back. Never `vi.mock` the `wailsjs/` modules — the
  contracts seam exists to avoid exactly that.
- The fake holds state rather than recording calls: drive it through the seam, then query
  `backend.state`. Emit events with `backend.emit(Event.NEW_LOG_ENTRY, {...})`. Drive a
  failure by replacing one method: `backend.data.removeCommand = async () => { throw ... }`.
- `installRecordingTerminals()` (from `@/testing/terminals.ts`) does the same for the
  emulator; use fake timers to drive the flush loop.
- Build objects with the builders in `@/testing/builders/`. They mutate and return `this`,
  so reuse one only when the earlier `with` calls should carry over.
- `installTranslations()` makes every key echo itself, so assertions name the key rather
  than the English copy. Assert toasts with `vi.spyOn(toast, "success" | "error")`.
- Mirror the Go conventions: Arrange / Act / Assert comments, and `sut` for the single
  unit under test.
- A hook owning an effect is driven through a probe component with `createRoot` and `act`
  (see `hooks/useAutosavedForm.test.tsx`). Anything a hook can hand to a plain function
  belongs in that function, tested without a renderer.
