# An interface earns its place by inverting a dependency, not by having a mock

Every use case exposes a single `Execute`, and for a long time every use case also
had an interface in front of it: `GetCommands`, `RunCommand`, `StopCommandGroup`
and the rest, each with exactly one implementation, `DefaultGetCommands` and its
siblings. The same held for four event handlers. #257, #283 and #288 removed them.
The registry now holds concrete pointers, and `eventbus.EventHandler` is the only
handler contract.

The rule that survived the removal is narrower than "delete single-implementor
interfaces", because plenty of the interfaces still here have one implementation:

- **A port that inverts a dependency across a layer boundary stays.** The
  repositories are declared in `domain/` and implemented in `infrastructure/`;
  `dialog.Dialogs` is declared where it is consumed and implemented by Wails; the
  facades front the filesystem, the OS and the toolkit. Each has one production
  implementation and each is load-bearing anyway, because the point is which way
  the arrow points, not how many types satisfy it. Deleting them would make the
  domain import GORM.
- **A contract with genuinely many implementers stays.** `eventbus.EventHandler`
  and `eventbus.Event` are satisfied by every handler and every event.
- **An interface whose only second implementation is a testify mock goes.** That
  is the whole of the deleted set.

The last case is the one worth writing down, because #283 got it wrong in a way
that reads as reasonable. It kept seven interfaces on the grounds that they "had a
second implementation" — and the second implementation was the mock written to
satisfy them. The justification is circular: the mock exists because the interface
exists, and the interface is then kept because the mock implements it. Nothing
outside the test suite ever needed the indirection.

The alternative we rejected is the familiar one: keep the interface so a consumer
can be tested against a mock. ADR-0004 already settled that backend behaviour is
verified through `apptest`, which wires the real use cases and fakes only what
leaves the process. A mock of a use case tests that a caller called it, which is a
restatement of the caller's source, not evidence about behaviour. When #288 removed
the interfaces the third-party HTTP tests depended on, the fix was to drive the
real backend through the harness — and two subtests that existed only to make a
mock return an error were dropped, because the harness has no way to break an
in-memory database and inventing one would have added a seam to prove a branch
nothing else needed.

## Consequences

Adding a use case is one type and one constructor, with no interface and no mock,
and the registry field is `*usecase`. A caller that wants to substitute it in a
test is doing something ADR-0004 says not to do.

The cost lands on anyone testing a consumer of a use case: there is no seam at that
point, so the test goes through `apptest` instead. That is the intended pressure.
When a genuine second implementation appears — a second runner, a second transport
— the interface comes back, declared in the package that consumes it.
