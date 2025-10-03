# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the frontend for Gomander, a Wails-based GUI application. The frontend is built with React 19, TypeScript, and communicates with the Go backend through Wails-generated bindings.

## Technology Stack

- **React 19** with TypeScript
- **Vite** as build tool with SWC for fast compilation
- **React Router 7** for routing (using HashRouter for Wails compatibility)
- **Zustand** for state management (vanilla stores + React hooks)
- **TailwindCSS v4** for styling
- **shadcn/ui** components (Radix UI primitives)
- **react-hook-form** + **Zod** for forms and validation
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
```

Note: Development server is run from the root via `make dev` or `wails dev` - not from this directory.

## Architecture

The frontend follows a **layered architecture** with clear separation of concerns:

### Directory Structure

```
src/
├── contracts/         # Wails backend interface layer (ONLY place wailsjs imports allowed)
│   ├── service.ts     # Exported service wrappers for backend calls
│   └── types.ts       # Type definitions from backend
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
│   ├── ui/            # shadcn components
│   ├── layout/        # Layout components
│   ├── modals/        # Modal dialogs
│   ├── inputs/        # Form inputs
│   └── utility/       # Utility components (EventListenersContainer, etc.)
├── hooks/             # Custom React hooks
├── contexts/          # React contexts (theme, version, etc.)
├── helpers/           # Pure utility functions
├── lib/               # Third-party library configurations
├── types/             # TypeScript type definitions
└── constants/         # Application constants
```

### Key Architectural Patterns

#### 1. Contracts Layer (Wails Abstraction)
- **CRITICAL**: Direct `wailsjs` imports are ONLY allowed in `src/contracts/`
- ESLint enforces this rule - imports from `wailsjs` anywhere else will fail linting
- All backend communication must go through `contracts/service.ts`
- This provides a clean abstraction layer between frontend and Wails-generated code

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
- Use `cn()` utility from `lib/utils.ts` to merge Tailwind classes
- shadcn components in `components/ui/` provide consistent design system

### Forms

- `react-hook-form` for form state management
- `zod` for schema validation
- `@hookform/resolvers` for integration

## Development Guidelines

### Import Organization
- ESLint enforces import sorting via `eslint-plugin-simple-import-sort`
- Imports are automatically sorted: external packages → internal imports (via @/)

### Adding New Backend Calls
1. Wails generates bindings in `wailsjs/` (DO NOT edit manually)
2. Add wrapper to `contracts/service.ts` (exported from `dataService`, `helpersService`, etc.)
3. Create use case in `useCases/` that calls the service
4. Never import from `wailsjs/` outside of `contracts/` directory

### Working with State
- In **components**: use `useCommandStore(state => state.field)` hooks
- In **use cases**: use `commandStore.getState()` or `commandStore.setState()`
- Keep selectors focused - only select what you need

### Adding Event Listeners
- All event listeners registered in `EventListenersContainer.tsx`
- Define event type in `contracts/types.ts` if new
- Call use case to handle state updates

## Known Constraints

- Must use `HashRouter` instead of `BrowserRouter` (Wails requirement)
- Cannot use interactive terminal commands in the app (no PTY support)
- Auto-generated `wailsjs/` code should never be manually edited
