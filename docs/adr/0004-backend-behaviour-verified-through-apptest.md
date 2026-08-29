# Backend behaviour is verified through apptest

`internal/apptest` wires the real use case registry, repositories, event bus and
database, and fakes only where the backend leaves the process: spawning
processes, the desktop runtime, and the filesystem. It is the default seam for
verifying backend behaviour — a new behaviour is tested through it, not by
stubbing a repository under a single use case.

The alternative is the mock-repository-per-use-case test, which is what the
codebase started with. Those tests pin the calls a use case makes rather than
what the user gets, so moving a rule between a use case, the domain and the
repository breaks them even when behaviour is identical. The recent refactors of
exactly that shape — the cascade rule moving into the domain, dropping
single-implementor interfaces — were cheap because the apptest suite never named
the collaborators.

## Consequences

Tests pay for a real SQLite database and a full migration run per harness
(`internal/testdb`), which is slower than a mock and means a schema mistake
surfaces as a test failure rather than in production. Both packages are excluded
from coverage — they are test infrastructure.

**This is a direction, not an accomplished state.** Around fifteen use case and
handler test files still drive `Mock*Repository` from the `domain/test` packages.
They are not wrong to exist and there is no migration in flight; they are simply
not the pattern to copy when writing a new test.
