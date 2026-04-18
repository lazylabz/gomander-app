# Contributing to Gomander

We welcome contributions to Gomander! This guide will help you get started with contributing to the project.

> [!IMPORTANT]
> Before starting work on an issue:
>  - Verify it's in the main project (if not, it probably means that the core maintainers haven't checked it yet)
>  - Comment on the issue to claim it, allowing maintainers to confirm it's prioritized and fully refined before you begin

## Getting Started

### Prerequisites
- Go 1.23+ installed
- Node.js and pnpm installed
- Follow the [Wails installation guide](https://wails.io/docs/gettingstarted/installation)

### Setup Steps
1. Fork the repository
2. Clone your fork locally
3. Follow the development setup instructions in the [README](README.md#development)
4. Create a new branch for your changes

### Project Structure
- `cmd/gomander/` - Main application entry point and Wails configuration
- `internal/` - Core business logic organized by domain (command, project, etc.)
- `cmd/gomander/frontend/` - React/TypeScript frontend
- `migrations/` - Database migrations (Go files using Goose)

## Development Guidelines

### Code Style

#### Backend (Go)
- Follow standard Go formatting with `gofmt`
- Use meaningful variable and function names
- Avoid writing unnecessary comments; code should be self-explanatory
- If some logic is complex, add concise comments to explain it
- Keep functions focused and small when possible
- Follow Clean Architecture principles: separate domain, application, and infrastructure layers
- For side effects, use domain events and event handlers
- Use dependency injection through the main.go registration process
- For database migrations, create `.go` files in `migrations/` using [goose](https://github.com/pressly/goose) format (see existing migration files for examples)
  - Use `goose create <name>` to create a new migration file

#### Frontend (TypeScript/React)
- Use TypeScript for type safety
- Follow React best practices and hooks patterns
- Use consistent naming conventions (camelCase for variables/functions, PascalCase for components)
- Keep components focused and reusable
- Leverage shadcn components for consistency
- Follow TailwindCSS conventions for styling

### Testing

Before submitting your changes:

1. **Run existing tests**: Ensure all current tests pass
2. **Add tests for new features**: Include unit tests for new functionality
3. **Test manually**: Run the application in dev mode and verify your changes work as expected
4. **Test builds**: Verify that the application builds successfully

```bash
# Run Go tests
go test ./...

# Run linting and formatting checks
make lint

# Run in development mode (from cmd/gomander directory)
cd cmd/gomander && wails dev

# Or use the Makefile shortcut from root
make dev

# Build the application
make all
```

> [!NOTE]
> At the time this document was written, we currently don't have frontend tests.

## Commit Guidelines

### Commit Message Format

Use [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) format:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

#### Types
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

#### Examples
```
feat(api): add endpoint for command group management
fix(ui): resolve command status display issue
docs: update installation instructions
```

### Commit Best Practices

- Make atomic commits (one logical change per commit)
- Write clear, descriptive commit messages
- Keep commits focused and small
- Avoid mixing unrelated changes in a single commit

## Pull Request Process

1. **Create a feature branch** from `main`
2. **Make your changes** following the guidelines above
3. **Test thoroughly** - both automated and manual testing
4. **Update documentation** if your changes affect user-facing functionality
5. **Submit a pull request** with:
   - Clear title describing the change
   - Detailed description of what was changed and why
   - Screenshots or examples if applicable
   - Reference to any related issues

### Pull Request Checklist

- [ ] Code follows project style guidelines
- [ ] Tests pass locally (`go test ./...`)
- [ ] Linting and formatting checks pass (`make lint`)
- [ ] Application builds successfully (`make all`)
- [ ] Documentation updated (if needed)
- [ ] Commit messages follow conventional format
- [ ] Changes are focused and atomic

## Reporting Issues

When reporting bugs or requesting features:

1. **Search existing issues** first to avoid duplicates
2. **Use the issue templates** if available
3. **Provide detailed information**:
   - Operating system and version
   - Steps to reproduce (for bugs)
   - Expected vs actual behavior
   - Screenshots or error logs if relevant

## Code Review

- All submissions require code review
- Be respectful and constructive in reviews
- Address feedback promptly
- Ask questions if feedback is unclear

## Adding a New Language

The backend auto-discovers available languages by reading the `cmd/gomander/locales/` directory at build time — there is no hardcoded list of supported locales anywhere in the codebase. Adding a new language requires creating one JSON file and registering a display label. No other backend or frontend code needs to change.

### Steps

#### 1. Create the locale file

Copy `cmd/gomander/locales/en.json` to `cmd/gomander/locales/<code>.json`, where `<code>` is a BCP-47 / ISO 639-1 language code (e.g. `fr`, `de`, `pt-BR`).

```bash
cp cmd/gomander/locales/en.json cmd/gomander/locales/fr.json
```

#### 2. Translate all values

Open the new file and translate every **value** (the right-hand side of each key-value pair). Do not rename or remove any keys.

Keep the following as-is:

- **Placeholders** like `{{version}}` or `{{count}}` — they are replaced at runtime
- **Plural suffixes** (`_one`, `_other`) — they must be present; some languages need additional forms (see the [i18next plurals guide](https://www.i18next.com/translation-function/plurals) for CLDR rules)
- **Inline HTML tags** — one key (`userSettingsForm.envPathsHelpBody`) contains `<code>` tags used by the `<Trans>` component; preserve them

Use `en.json` as the canonical reference.

#### 3. Register the language label

Add an entry to `cmd/gomander/frontend/src/constants/languages.ts`:

```ts
export const LANGUAGE_LABELS: Record<string, string> = {
  en: 'English',
  es: 'Español',
  fr: 'Français', // add your language here
};
```

The label must be the language's **native name** (as it appears in that language), because the selector shows it regardless of the active locale.

#### 4. Rebuild the application

The locale files are compiled into the binary via `go:embed`. A rebuild is required for the new file to be included:

```bash
# Development
make dev

# Production build
make all
```

#### 5. Verify manually

1. Open **Settings → Language** — your language should appear in the dropdown
2. Switch to it and verify the full UI: main screens, modals, form validation error messages, and toast notifications
3. Check that no strings appear empty (would indicate a missing key)

### Commit format

```
feat(i18n): add <language name> translations
```

## Getting Help

- Check existing documentation and issues first
- Open a discussion for questions about contributing
- Reach out to maintainers if you need guidance

Thank you for contributing to Gomander! 🚀
