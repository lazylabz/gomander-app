# CLAUDE.md

The Gomander frontend: React 19, TypeScript, Vite, Zustand, TailwindCSS v4,
shadcn/ui, react-hook-form + Zod, Biome, pnpm. It talks to the Go backend through
Wails bindings. See [the root CLAUDE.md](../../../../.claude/CLAUDE.md) for the backend.

## Where the reasoning lives

- **[CONTEXT.md](../../../../CONTEXT.md)** — the glossary, shared with the backend.
- **[docs/adr/](../../../../docs/adr/)** — decisions and the alternatives they rejected.
  ADR-0001 (ports at the third-party edges), 0002 (autosave), 0003 (browser storage),
  0005 (coverage) and 0008 (component tests) all bind code in this directory.

`pnpm run` lists the scripts. The dev server runs from the repo root (`make dev`).

## Layout

```
src/
├── contracts/         # the backend seam - the ONLY place the Wails bindings are imported
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
case. Biome fails any `bindings/` or `@wailsio/runtime` import outside `contracts/`, and
any `@xterm/*` import outside `commandOutput/`.

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

Vitest, happy-dom, Testing Library, tests next to the code as `*.test.ts(x)`. ADR-0008
sets the strategy: a rule is tested in the deepest module that owns it, and a component
is rendered only for its wiring. Logic found in a component moves out before it is
tested.

- `installInMemoryBackend()` (from `@/testing/backend.ts`) swaps every service;
  `resetBackendServices()` puts Wails back. Never `vi.mock` the `bindings/` modules — the
  contracts seam exists to avoid exactly that.
- The fake holds state rather than recording calls: drive it through the seam, then query
  `backend.state`. Emit events with `backend.emit(Event.NEW_LOG_ENTRY, {...})`. Drive a
  failure by replacing one method: `backend.data.removeCommand = async () => { throw ... }`.
- `installRecordingTerminals()` (from `@/testing/terminals.ts`) does the same for the
  emulator; use fake timers to drive the flush loop.
- Build objects with the builders in `@/testing/builders/`. They mutate and return `this`,
  so reuse one only when the earlier `with` calls should carry over.
- `installTranslations()` makes every key echo itself, so assertions name the key rather
  than the English copy. A test that must see an interpolated value gives that one key
  copy with `withTranslation(key, copy)`. Assert toasts with `vi.spyOn(toast, "success" | "error")`.
- Mirror the Go conventions: Arrange / Act / Assert comments, and `sut` for the single
  unit under test.
- Render components with `renderWithProviders` (from `@/testing/render.tsx`): it wraps
  the router and the sidebar provider, returns a user-event `user`, and `location()` for
  asserting navigation. Query by role and label; assert with the jest-dom matchers.
- Never `vi.mock` a child component or a use case: render the real ones against the
  in-memory backend and assert on what the user sees or on `backend.state`.
- Stores are module singletons: call `resetStores()` (from `@/testing/stores.ts`) in
  `beforeEach` instead of seeding every field by hand. A new store is added there.
- A hook owning an effect is driven through a probe component (see
  `hooks/useAutosavedForm.test.tsx`). Anything a hook can hand to a plain function
  belongs in that function, tested without a renderer.
- A test rendering `App` (see `App.test.tsx`) calls `cleanup()` before
  `resetBackendServices()` in its `afterEach`. The global cleanup in
  `testing/setup.ts` runs after the file's own hooks, and `EventListenersContainer`
  unsubscribes through the backend on unmount, which by then is the Wails stub that
  throws. `App` also runs the real `initI18n`, so the in-memory backend needs
  `translations: { en: ... }` or the App stays on its load-failed screen.
- Never click a dnd-kit drag handle. The sidebar's `DndContext` sets no activation
  constraint, so a pointerdown starts a drag that puts a capture-phase click blocker
  on the document, and dnd-kit removes it 50 ms after release. The test passes and
  the next tests in the file lose their clicks. Test drag rules in the module that
  owns them (ADR-0008).

## Verify in a browser

`pnpm dev` serves the frontend alone on 127.0.0.1:9245. That port is not the Wails
runtime: every backend call is a `POST /wails/runtime`, which 404s there, so the
page opens on "Something went wrong while loading". Only the native window
(`make dev`) reaches the Go backend.

To see the UI in a browser, install the in-memory backend from the page and press
Retry; the app's modules are the ones the page imports:

```js
const { installInMemoryBackend } = await import("/src/testing/backend.ts");
const { ProjectBuilder } = await import("/src/testing/builders/project.ts");
const { CommandBuilder } = await import("/src/testing/builders/command.ts");
const project = new ProjectBuilder().withId("p1").withName("Demo").build();
window.backend = installInMemoryBackend({
	projects: [project],
	currentProject: project,
	commands: [new CommandBuilder().withProjectId("p1").withName("api").build()],
	translations: { en: {} },
});
```

A key missing from `translations` renders as the key itself. For the English copy,
inline the contents of `cmd/gomander/locales/en.json`: it sits outside what Vite
serves, so the page cannot import it. Show a backend failure the same way a test
does, by replacing one method (`backend.data.getUserConfig = async () => { throw
new Error() }`) before Retry.
