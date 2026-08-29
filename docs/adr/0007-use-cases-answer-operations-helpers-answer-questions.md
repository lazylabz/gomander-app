# Use cases answer operations, helpers answer questions

The frontend reaches the backend through two kinds of object bound to Wails. One is
`WailsControllers`, which owns nothing but a use case registry and calls `Execute`.
The other is a handful of helpers — `UiPathHelper`, `UIOsHelper`, `UIFsHelper` —
bound directly, each holding its own collaborators.

Nothing said which was which, so the update flow ended up on the helper side: a
`ReleaseHelper` that fetched an Atom feed, downloaded a binary, launched it and quit
the app. It could not be driven from `apptest`, the third-party HTTP server could not
reach it, and it had to be handed a `context.Context` after construction because the
shell only produces one at startup. #300 moved it. The rule that decides where the
next call goes:

- **It is a use case when it is an operation the app performs**: it reads or changes
  what the app holds, reaches the network or the filesystem on its own initiative, or
  is a step in a flow the backend owns. It goes in the registry, is reached through a
  controller that adds nothing, and is drivable from `apptest` and from the
  third-party HTTP server.
- **It is a helper when it is a question the UI needs answered to render or to fill
  in a form, and answering it leaves the app exactly as it was**: which OS this is,
  what a working directory resolves to, which directory the user just picked, showing
  a folder in the file manager. A helper holds no application state and never appears
  in the registry.
- **When both readings fit, the flow wins.** Asking which Release is running is a
  question by itself, but it is the first step of checking for an update, so it sits
  with the other three rather than splitting the flow across two architectures.

The frontend contract mirrors the split: `DataService` in `contracts/ports.ts` is the
use cases, `HelpersService` is the helpers.

## Alternatives rejected

**Everything becomes a use case.** It reads as the tidier rule, and it is what the
Clean Architecture layout suggests. But `GetComputedPath` is a domain function the UI
calls to preview a path as the user types, and `GetOs` is a build constant; wrapping
them in a registry field, a controller and an `apptest` fake buys nothing, and the
third-party API has no use for either. The pressure would be to route UI concerns
through the application layer to satisfy a rule rather than a need.

**Keep binding whatever is convenient.** This is the status quo the update flow came
from. The cost is not hypothetical: three operations sat outside every seam the rest
of the backend is verified through, and one object existed half-built waiting for the
shell to finish it.

## Consequences

Adding an operation is the four frontend steps of ADR-0001 plus a registry field, a
constructor call in `buildDeps` and a controller method — and, because
`TestTheHarnessWiresEveryOperationTheAppExposes` fails on a registry field the harness
leaves nil, wiring it into `apptest` too.

The three helpers left are the whole list. A fourth has to argue from this rule, and
"it is easier to bind than to add to the registry" is not that argument.
