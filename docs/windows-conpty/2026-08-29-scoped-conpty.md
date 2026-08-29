# Scoped Windows ConPTY Implementation Plan

> **For agentic workers:** Execute these steps in order and verify each behavior before moving on.

**Goal:** Use ConPTY only on configured, verified Windows host environments, initially Windows 10 workstations with ConPTY support, while preserving the existing pipe runner everywhere else.

**Architecture:** Application-level static parameters live at the composition root and inject a `runner.Config` into `DefaultRunner`. The Windows runner resolves the current host environment from `RtlGetVersion`, chooses ConPTY only when that environment appears in the configured allowlist and `CreatePseudoConsole` exists, and otherwise uses the existing `os.Pipe` execution path. Unix behavior remains the existing pipe behavior.

**Tech Stack:** Go 1.23+, Wails v2, `golang.org/x/sys/windows`, `github.com/aymanbagabas/go-pty`.

---

### Task 1: Centralize application parameters

**Files:**
- Create: `cmd/gomander/application_config.go`
- Modify: `cmd/gomander/main.go`
- Create: `internal/runner/config.go`

- [x] Define `runner.Config` with an explicit ConPTY host-environment allowlist and the `windows10` environment identifier.
- [x] Define the application configuration at the composition root with `windows10` as the only ConPTY-enabled environment.
- [x] Move the existing configuration-directory name into the same application configuration and inject the runner configuration into `NewDefaultRunner`.

### Task 2: Preserve pipe execution behind a process contract

**Files:**
- Create: `internal/runner/process.go`
- Create: `internal/runner/process_exec.go`
- Modify: `internal/runner/proc_unix.go`
- Modify: `internal/runner/runner.go`
- Test: `internal/runner/process_exec_test.go`

- [x] Add the smallest process contract needed by the shared runner: start, wait, PID, and ConPTY-only output filtering.
- [x] Move the current explicit `os.Pipe` ownership into `execProcess` without changing its wait, drain, exit-error, or environment behavior.
- [x] Keep Unix command construction and process-group shutdown behavior unchanged.
- [x] Verify trailing output and non-zero exit errors remain observable through the pipe path.

### Task 3: Gate ConPTY to configured Windows 10 hosts

**Files:**
- Modify: `internal/runner/proc_windows.go`
- Test: `internal/runner/proc_windows_test.go`

- [x] Resolve Windows 10 workstation builds with `RtlGetVersion`: build 17763 through 21999 and `VER_NT_WORKSTATION` only.
- [x] Select ConPTY only when `windows10` is configured and the `CreatePseudoConsole` API exists.
- [x] Preserve the existing hidden pipe process on every non-enabled host.
- [x] Port ConPTY start, output cleanup, exit status, and shutdown behavior from PR #227 without enabling it outside the allowlist.
- [x] Cover pre-1809 Windows 10, supported Windows 10, Windows Server, and Windows 11 classification.

### Task 4: Verify behavior and scope

**Files:**
- Test: `internal/runner/runner_test.go`
- Test fixture: existing external `uvicorn_repro_test.go` in the retained experiment environment

- [x] Run `gofmt` on changed Go files.
- [x] Run targeted runner tests on Windows.
- [x] Run `go test ./... -count=1` because the runner constructor and shared execution infrastructure change.
- [x] Run `go vet ./...`.
- [x] Compile the runner tests for Linux amd64.
- [x] Run the retained Uvicorn 0.50.0/watchfiles 1.2.0 reproduction with the Windows 10 allowlist enabled and confirm three reloads with and without keep-alive.
- [x] Review the diff against `origin/main` and confirm Windows 11, Windows Server, old Windows 10, Unix behavior, persistence, and frontend configuration remain unchanged; the runner constructor changes only to accept the injected compatibility configuration.
