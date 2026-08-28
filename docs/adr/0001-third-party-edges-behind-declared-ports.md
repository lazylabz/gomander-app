# Third-party edges live behind a declared port

The frontend touches two libraries it does not control: the Wails bindings the Go
backend generates, and the xterm.js terminal emulator. Both are confined to one
directory each — `src/contracts/` and `src/commandOutput/` — with the boundary
enforced by Biome's `noRestrictedImports`, so a stray import fails linting rather
than review.

The contract itself is **declared** in `contracts/ports.ts` as hand-written types,
and the Wails adapter is checked against it with `satisfies`. The obvious
alternative is to derive the contract from the generated bindings: `typeof
import("wailsjs/...")`. We rejected it. Deriving makes the generator the author of
our contract — a regenerated binding silently reshapes it, and the type error
lands wherever the app happens to call it instead of at the boundary. Declaring it
means a backend change that breaks the contract fails in exactly one file.

## Consequences

Every new backend call is four steps, not one: declare it in `ports.ts`, map it in
`adapters/wails.ts`, implement it in `adapters/inMemory.ts`, then use it. That
cost is the point — the in-memory adapter is what lets tests run with no Wails
runtime and no `vi.mock` of generated code.
