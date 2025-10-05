
# Copilot Instructions for AI Coding Agents

## Project Overview
- This is an Android-focused fork of the PC-based Arduino CLI. It contains custom modifications and may not behave like the upstream project.
- Major components:
  - `main.go`: Entry point for CLI logic and direct CLI interactions.
  - `commands/`: Service implementations for board, library, platform, sketch, and upload management. Each `service_*` file is a distinct CLI feature.
  - `internal/`: Core logic, device communication, and integration test support.
  - `mobile/`: Android-specific interfaces and relay functions, including file type conversions for mobile use.
  - `mobile/bridge/` (planned): Will contain bridge functions to allow Go code to communicate with the Android app and request serial info. Functions are relayed in `mobile/mobile.go`.
  - `docs/`: Extensive documentation for architecture, integration, and platform specs.

## Build & Test Workflows
- **Standard build:**
  - Run `task build` (uses Go modules, see `Taskfile.yml`).
  - For Android/mobile builds on Windows, use `build_arduinobuddy.bat`.
- **Unit tests:**
  - Run `task go:test` for Go-based tests.
- **Integration tests:**
  - Requires attached Arduino hardware (see `docs/CONTRIBUTING.md`).
  - Run `task go:integration-test`.
  - To run a specific package: `go test -v github.com/arduino/arduino-cli/internal/integrationtest/lib`
  - To run a specific test: `go test -v ... -run TestLibUpgradeCommand`
- **Linting/Formatting:**
  - Run `task check` for style and lint checks.
  - Run `task go:format` for Go code formatting.
  - Run `task general:format-prettier` for Prettier formatting.

## Key Architectural Patterns
- **Service boundaries:** Each CLI feature is implemented as a service in `commands/`, with clear separation (e.g., board, library, platform).
- **Android integration:** The `mobile/` directory provides interfaces and relays for Android, converting file types and handling communication. Planned `mobile/bridge/` will enable Go-to-Android serial info requests.
- **gRPC integration:** See `docs/integration-options.md` for gRPC API usage and embedding guidance. Prefer gRPC for cross-process communication.
- **Platform configuration:** Platform specs and board definitions use key-value config files (see `docs/platform-specification.md`). Properties can reference other properties using `{property}` syntax.
- **Testing:** Integration tests run the CLI in a separate process and validate real device interactions.

## Project-Specific Conventions
- Android and Windows are the primary supported platforms for custom features (see `README.md`).
- Use `Taskfile.yml` for all build/test/lint workflows; do not rely on raw Go commands unless customizing.
- Board/platform config files (e.g., `boards.txt`) use custom debug and build flags.
- When modifying gRPC APIs, see `docs/UPGRADING.md` for protocol changes and migration notes.

## External Dependencies
- Go 1.21+ required.
- Uses `buf` for protobuf/gRPC codegen (see `docs/CONTRIBUTING.md`).
- Linting requires `golangci-lint` and `licensed` (see `Taskfile.yml`).

## Examples
- To add a new CLI command, create a new `service_*` file in `commands/` and register it in `main.go`.
- To update platform specs, edit config files in `docs/platform-specification.md` and related directories.
- To extend Android integration, add relay or bridge functions in `mobile/` and (when created) `mobile/bridge/`.

---
For unclear or incomplete sections, please provide feedback so this guide can be improved for future AI agents.
