# Frontend components are tested for their wiring, not their logic

A rule lives in the deepest module that can own it, and is tested there: a plain
function with no renderer, or a use case, query or store driven through the
in-memory backend. A component is rendered only for what a render alone can prove —
which control shows for which state, that a click reaches the right use case, that a
modal closes on success and stays open with the user's input on failure.

So before a component is tested, the logic inside it moves out. The drag rules of
`CommandGroupCommandsField` became `commandGroupSelection.ts`, the group rule of the
import modal became `blueprintSelection.ts`, `useShortcut` gave its matching to
`matchesShortcut`, and the release flow left its React context for a store and three
use cases. dnd-kit cannot be dragged in a DOM emulator at all, so for those rules the
extraction was the only way to test them.

## Rejected alternatives

- **Rendering everything and asserting through the DOM.** It works for a while, but
  every rule then needs a full provider tree and a scripted interaction to reach, and
  a rule that moves between a component and a use case breaks its tests even when the
  behaviour holds — the same trap ADR-0004 describes for mock repositories.
- **Mocking child components or use cases with `vi.mock`.** A mocked use case proves
  only that the component calls it, which restates the source. Components render
  their real children and reach the real use cases; the fakes stay at the edges of
  ADR-0001, the backend and the terminal emulator.
- **jsdom.** The suite ran on it before any component was rendered. A spike with the
  Radix primitives the app uses showed `Select` failing under user-event
  (`hasPointerCapture is not a function`) without polyfills, while happy-dom passed
  them all and ran the existing suite in under half the time.

## Consequences

Testing Library, user-event and jest-dom matchers are the tools for a render;
`renderWithProviders` in `src/testing/` supplies the router and sidebar providers and
exposes the current location. Stores are module singletons, so tests reset them with
`resetStores()` rather than by hand.

`design-system/` (vendored shadcn), the Wails and xterm adapters and presentational
components with no branching stay untested on purpose.
