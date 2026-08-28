# CLAUDE.md

Gomander is a cross-platform desktop app (Wails: Go + React) for launching,
monitoring and organising the shell commands of a project.

## Where the reasoning lives

- **[CONTEXT.md](../CONTEXT.md)** — the glossary. It is the naming authority: use its
  terms in code, tests and commits, and add a term there when a new one settles.
- **[docs/adr/](../docs/adr/)** — decisions and the alternatives they rejected. Read the
  relevant one before "fixing" something that looks wrong on purpose.
- **[frontend](../cmd/gomander/frontend/.claude/CLAUDE.md)** — everything under
  `cmd/gomander/frontend/`.
- Package docs carry the why for their own package — `openedproject`, `execution`,
  `ordering`, `commandgroup/domain`, `apptest`, `domainerrors`. A new rationale goes
  next to the code it explains, not into this file.

## Stack

Go 1.23+, Wails v2, SQLite through GORM, Goose migrations as `.go` files in
`migrations/`. `make help` lists the dev, test and build targets.

## Layout

```
cmd/gomander/
├── main.go              # buildDeps: every dependency is constructed and wired here
├── controllers.go       # what the Wails frontend can call
├── frontend/
└── thirdpartyserver/    # REST API on ports 9001-9100; spec in openapi.yaml

internal/
├── app/                 # startup, shutdown, event handler registration
├── usecases/            # the registry the UI and the third-party API reach
├── command/             # ─┐
├── commandgroup/        #  ├ one per entity: domain/ application/ infrastructure/
├── project/             #  │
├── config/              # ─┘ user configuration
├── openedproject/       # which Project is open
├── execution/           # the environment a Command runs in, and the Runner port
├── ordering/            # dense Positions within a list
├── unitofwork/          # one transaction spanning several repositories
├── runner/              # spawns Commands, streams their output
├── event/ eventbus/     # domain events and the in-memory bus
├── domainerrors/        # the errors every entity domain reports
├── localization/        # translations served to the frontend
├── releases/            # update checks and platform binaries
├── facade/              # fs, io, os, open and Wails runtime, so tests can fake them
├── dialog/              # asking the user for a path, and the Wails adapter that does
├── uihelpers/ helpers/  # what the frontend calls directly; array utilities
├── apptest/             # the seam backend behaviour is verified through
└── testdb/              # in-memory SQLite plus migrations

migrations/              # goose create <name>, Go files with Up and Down
```

## Architecture rules

- Clean Architecture: domain → application → infrastructure. The domain layer has no
  external dependencies.
- One use case per operation, exposing a single `Execute`.
- Side effects travel as domain events on the event bus and land in
  `application/handlers`, never inline in the use case that caused them.
- Repository interfaces are declared in `domain/` and implemented in `infrastructure/`.
- Construct and wire every dependency in `buildDeps`.

## Conventions

- Conventional commits: `feat(scope):`, `fix`, `docs`, `refactor`, `test`, `chore`.
- Test new backend behaviour through `internal/apptest` (see ADR-0004): arrange with
  the domain builders, act through `h.UseCases`, assert on what the user or the OS
  sees — `StartedProcesses`, `EmittedEvents`, a read use case.
- Run `go test ./...` and `make lint` before handing work back.

## Limitations

- No PTY: the runner reads stdout, so TUI commands (ngrok and the like) do not work.
- macOS builds are unsigned — they need `sudo xattr -d com.apple.quarantine`.
