# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the frontend for Gomander, a Wails-based GUI application. The frontend is built with React 19, TypeScript, and
communicates with the Go backend through Wails-generated bindings.

## Technology Stack

- **React 19** with TypeScript
- **Vite** as build tool with SWC for fast compilation
- **React Router 7** for routing (using HashRouter for Wails compatibility)
- **Zustand** for state management (vanilla stores + React hooks)
- **TailwindCSS v4** for styling
- **shadcn/ui** components (Radix UI primitives)
- **react-hook-form** + **Zod** for forms and validation
- **Biome** for linting and formatting
- **pnpm** as package manager

## Common Commands

### Development

```bash
# Type checking (most commonly used during development)
pnpm run typecheck

# Linting
pnpm run lint

# Auto-fix lint issues
pnpm run lint:fix

# Auto-fix lint issues, including fixes that Biome classifies as unsafe
pnpm run lint:fix:unsafe

# Tests (Vitest)
pnpm run test

# Tests in watch mode
pnpm run test:watch

# Tests with coverage
pnpm run test:cov

# A single test file
pnpm run test src/queries/fetchCommands.test.ts
```

Note: Development server is run from the root via `make dev` or `wails dev` - not from this directory.

## Architecture

The frontend follows a **layered architecture** with clear separation of concerns:

### Directory Structure

```
src/
├── contracts/         # Wails backend interface layer (ONLY place wailsjs imports allowed)
│   ├── ports.ts       # The backend contract, declared as types
│   ├── service.ts     # The services the app calls, plus the adapter swap for tests
│   ├── types.ts       # Type definitions from backend
│   └── adapters/      # Implementations of the contract
│       ├── wails.ts   # Generated bindings (production)
│       └── inMemory.ts # Fake holding queryable state (tests)
├── useCases/          # Business logic organized by domain
│   ├── command/       # Command-specific operations (start, stop, create, etc.)
│   ├── commandGroup/  # Command group operations
│   ├── project/       # Project operations
│   ├── userConfig/    # User configuration
│   └── logging/       # Logging operations
├── queries/           # Data fetching operations (read-only)
├── store/             # Zustand state stores (one per domain)
│   ├── commandStore.ts
│   ├── commandGroupStore.ts
│   ├── projectStore.ts
│   └── userConfigurationStore.ts
├── screens/           # Top-level screen components
│   ├── ProjectSelectionScreen/
│   ├── SettingsScreen/
│   └── LogsScreen/
├── components/
│   ├── layout/        # Layout components
│   ├── modals/        # Modal dialogs
│   ├── inputs/        # Form inputs
│   └── utility/       # Utility components (EventListenersContainer, etc.)
├── design-system/     # Design system (shadcn/ui components and utilities)
│   ├── components/
│   │   └── ui/        # shadcn/ui components (Button, Dialog, etc.)
│   ├── hooks/         # shadcn/ui hooks (use-mobile, etc.)
│   └── lib/           # shadcn/ui utilities (cn, utils, etc.)
├── hooks/             # Custom React hooks (application-specific)
├── contexts/          # React contexts (theme, version, etc.)
├── helpers/           # Pure utility functions
├── testing/           # Test-only helpers (in-memory backend installer, builders)
├── types/             # TypeScript type definitions
└── constants/         # Application constants
```

### Key Architectural Patterns

#### 1. Contracts Layer (Wails Abstraction)

- **CRITICAL**: Direct `wailsjs` imports are ONLY allowed in `src/contracts/`
- Biome enforces this rule - imports from `wailsjs` anywhere else will fail linting
- All backend communication must go through `contracts/service.ts`
- The contract is declared in `contracts/ports.ts` as types; adapters implement it and are
  checked with `satisfies` - the shape is not derived from the generated bindings
- Two adapters exist: `adapters/wails.ts` (production) and `adapters/inMemory.ts` (tests)

#### 2. State Management (Zustand)

- **Vanilla stores** created with `createStore` for use in non-React code (use cases)
- **React hooks** created with `useStore` for use in components
- Each domain has its own store (command, commandGroup, project, userConfig)
- Stores are kept minimal - just state and setters

Example pattern:

```typescript
// In store file
export const commandStore = createStore<CommandStore>()(...)  // For use cases
export const useCommandStore = <T>(selector) => useStore(commandStore, selector)  // For components
```

#### 3. Use Cases Pattern

- Business logic lives in `useCases/` directory, organized by domain
- Each use case is a single exported async function
- Use cases interact with backend via `contracts/service.ts`
- Use cases update state via vanilla Zustand stores (not React hooks)

Example: `useCases/command/startCommand.ts` calls backend and updates state

#### 4. Queries vs Use Cases

- **Queries** (`queries/`): Read-only data fetching operations
- **Use Cases** (`useCases/`): Write operations or complex business logic
- Queries typically load data into stores on app initialization

#### 5. Event System

- Backend pushes real-time updates via Wails events
- `EventListenersContainer` component (rendered in App.tsx) listens to all events
- Events defined in `contracts/types.ts` (Event enum)
- Event handlers call use cases to update state (e.g., PROCESS_STARTED → updateCommandStatus)
- Log buffering: Logs are buffered for 30ms before being flushed to state for performance

### Path Aliases

- `@/` maps to `src/` directory
- Configured in both `vite.config.ts` and `tsconfig.json`

### Routing

- Uses React Router 7 with `HashRouter` (required for Wails)
- Routes defined in `src/routes.ts`
- Main routes: ProjectSelection, Logs, Settings

### Styling

- TailwindCSS v4 configured via `@tailwindcss/vite` plugin
- Use `cn()` utility from `design-system/lib/utils.ts` to merge Tailwind classes
- shadcn components in `design-system/components/ui/` provide consistent design system

### Forms

- `react-hook-form` for form state management
- `zod` for schema validation
- `@hookform/resolvers` for integration

## Development Guidelines

### Import Organization

- Biome enforces import sorting via its `organizeImports` assist action
- Imports are automatically sorted: external packages → internal imports (via @/)

### Control Flow

- Biome enforces `style/useBlockStatements` - `if`/`for`/`while` bodies always need braces, even for a single statement
- Its autofix is classed unsafe, so `pnpm run lint:fix` skips it silently; use `pnpm run lint:fix:unsafe`

### Adding New Backend Calls

1. Wails generates bindings in `wailsjs/` (DO NOT edit manually)
2. Declare the method on the matching type in `contracts/ports.ts`
3. Map it in `contracts/adapters/wails.ts` and implement it in `contracts/adapters/inMemory.ts`
4. Create use case in `useCases/` that calls the service
5. Never import from `wailsjs/` outside of `contracts/` directory

### Working with State

- In **components**: use `useCommandStore(state => state.field)` hooks
- In **use cases**: use `commandStore.getState()` or `commandStore.setState()`
- Keep selectors focused - only select what you need

### Adding Event Listeners

- All event listeners registered in `EventListenersContainer.tsx`
- Define event type in `contracts/types.ts` if new
- Call use case to handle state updates

### Testing

Vitest, jsdom environment, tests live next to the code as `*.test.ts(x)`.

Node 22.10 or newer is required (`engines` in `package.json`): jsdom 30 pulls an undici that
calls `worker_threads.markAsUncloneable`. On an older Node every test file fails to start.

- Tests substitute the in-memory adapter at the seam: `installInMemoryBackend()` (from
  `@/testing/backend.ts`) swaps every service exported by `contracts/service.ts`, and
  `resetBackendServices()` puts the Wails adapter back
- The swap is module-level (the service exports are live `let` bindings) rather than
  injected into every use case: it reaches all call sites without changing a signature.
  `setBackendServices` is for tests only - production code never calls it
- The fake **holds state** instead of asserting interactions: drive it through the seam and
  then query `backend.state`, so refactors behind the contract do not break the tests
- Emit backend events on demand with `backend.emit(Event.NEW_LOG_ENTRY, { id, line })`
- Build domain objects with the builder classes in `@/testing/builders/`, mirroring the Go
  side: `new CommandBuilder().withId("cmd-1").build()`. They mutate and return `this`, so
  reuse a builder only when you mean the earlier `with` calls to carry over
- Mirror the Go test conventions: Arrange / Act / Assert comments, and name the unit under
  test `sut` when the test has a single one
- Never `vi.mock` the `wailsjs/` modules - that is what the contracts seam exists to avoid
- Frontend coverage is not uploaded to Codecov yet (it would fail the repo-wide 80% project
  target while the suite is this small)

## Known Constraints

- Must use `HashRouter` instead of `BrowserRouter` (Wails requirement)
- Cannot use interactive terminal commands in the app (no PTY support)
- Auto-generated `wailsjs/` code should never be manually edited
