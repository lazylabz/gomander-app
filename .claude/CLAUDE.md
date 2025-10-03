# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gomander is a cross-platform GUI application built with Wails (Go + React) for launching, monitoring, and organizing shell commands. The app allows users to organize commands by project, bundle them into groups, and manage their execution.

## Technology Stack

- **Backend**: Go 1.23+, Wails v2
- **Frontend**: React 19, TypeScript
    - @../cmd/gomander/frontend/.claude/CLAUDE.md
- **Database**: SQLite (via GORM)
- **Migrations**: Goose (Go migration files in `migrations/`)

## Common Commands

### Development
```bash
# Start development server (from root)
make dev

# Or directly with Wails (from cmd/gomander)
cd cmd/gomander && wails dev
```

### Testing and Quality
```bash
# Run all Go tests
go test ./...

# Run linting and type checking (backend + frontend)
make lint
```

### Building
```bash
# Build for all platforms
make all

# Build for specific platforms
make darwin      # macOS (both architectures + DMG)
make windows     # Windows (both architectures + installers)
make linux       # Linux (both architectures, uses Docker)

# See all build options
make help
```

### Single Test Execution
```bash
# Run a specific test
go test ./internal/command/domain -run TestCommandName

# Run tests in a specific package with verbose output
go test -v ./internal/command/application/usecases
```

## Architecture

Gomander follows **Clean Architecture** principles with clear separation of concerns:

### Directory Structure

```
cmd/gomander/
├── main.go              # Application entry point, dependency injection
├── controllers.go       # Wails frontend-backend controllers
├── frontend/            # React TypeScript application
└── thirdpartyserver/    # REST API for third-party integrations (ports 9001-9100)
    └── openapi.yaml

internal/                # Core business logic organized by domain
├── app/                 # Application lifecycle and dependency coordination
├── command/             # Command domain (run/stop individual commands)
│   ├── domain/          # Entities, value objects, repository interfaces
│   ├── application/     # Use cases and event handlers
│   └── infrastructure/  # GORM repository implementations
├── commandgroup/        # Command group domain (bundle and run multiple commands)
├── project/             # Project domain (organize commands by project)
├── config/              # User configuration domain
├── event/               # Domain events system
├── eventbus/            # In-memory event bus for domain events
├── runner/              # Command execution engine
├── logger/              # Logging abstraction
├── facade/              # Facades for OS operations (fs, runtime)
└── uihelpers/           # UI-specific helpers

migrations/              # Database migrations (Goose format .go files)
```

### Clean Architecture Layers

1. **Domain Layer** (`domain/`): Entities, value objects, repository interfaces, domain events. No external dependencies.

2. **Application Layer** (`application/`):
   - **Use Cases**: Business logic implementation (one use case per operation)
   - **Event Handlers**: React to domain events for side effects

3. **Infrastructure Layer** (`infrastructure/`): GORM repository implementations, external service integrations

### Key Architectural Patterns

- **Dependency Injection**: All dependencies are registered in `cmd/gomander/main.go` (`registerDeps` function) and injected through constructors
- **Domain Events**: Side effects are handled via domain events published through the event bus
- **Repository Pattern**: Data access abstracted behind repository interfaces defined in domain layer
- **Use Case Pattern**: Each business operation is a separate use case struct with a single `Execute` method

### Database

- SQLite database stored in user config directory (`~/Library/Application Support/gomander/data.db` on macOS)
- Migrations managed by Goose - create new migrations with `goose create <name>` in `migrations/` directory
- Migration files are Go files with `Up` and `Down` functions

### Third-Party API

The application exposes a REST API on ports 9001-9100 for third-party integrations. API spec: `cmd/gomander/thirdpartyserver/openapi.yaml`

## Frontend

For frontend-specific guidance, architecture, and development guidelines, see `cmd/gomander/frontend/CLAUDE.md`.

## Development Guidelines

### Backend

- Follow standard Go formatting with `gofmt`
- Maintain Clean Architecture: domain → application → infrastructure
- Use dependency injection through `main.go` registration
- For side effects, use domain events and event handlers
- Keep functions focused and small
- Code should be self-explanatory; add comments only for complex logic
- Create database migrations as `.go` files using `goose create <name>`

### Commit Messages

Use [conventional commits](https://www.conventionalcommits.org/):
- `feat(scope): description` - New feature
- `fix(scope): description` - Bug fix
- `docs: description` - Documentation changes
- `refactor(scope): description` - Code refactoring
- `test: description` - Test additions/corrections
- `chore: description` - Build/tooling changes

### Testing

- Run existing tests before submitting changes: `go test ./...`
- Add unit tests for new functionality
- Test manually in dev mode: `make dev`
- Verify builds succeed: `make all`

## Known Limitations

- TUI commands (e.g., ngrok) are not supported - the runner relies on stdout and doesn't support PTY
- macOS builds require manual quarantine removal (`sudo xattr -d com.apple.quarantine`) due to lack of code signing
