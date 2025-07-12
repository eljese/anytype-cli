# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the Anytype CLI, a Go-based command-line interface for interacting with Anytype. It uses a client-server architecture where the CLI communicates with a separate middleware server (anytype-heart) via gRPC.

## Build Commands

```bash
# Download the middleware server (required before first use)
make download-server

# Build the CLI
make build

# Install system-wide
make install

# Install user-local (~/.local/bin)
make install-local

# Manual download of middleware server
make download-server

# Run linting
make lint

# Manual build with version info
go build -ldflags "-X main.Version=$(git describe --tags --always) -X main.Commit=$(git rev-parse HEAD) -X main.BuildTime=$(date -u +%Y%m%d-%H%M%S)" -o dist/anytype
```

## Development Workflow

1. **Initial Setup**:
   ```bash
   make download-server  # Download anytype-heart middleware
   make build           # Build the CLI
   ```

2. **Running the Application**:
   ```bash
   # Start the daemon (required)
   ./dist/anytype daemon
   
   # In another terminal, start the server
   ./dist/anytype server start
   ```

3. **Code Formatting and Linting**:
   ```bash
   go fmt ./...
   go vet ./...
   make lint  # Uses golangci-lint
   ```

## Architecture Overview

### Command Structure (`/cmd/`)
- Uses Cobra framework for CLI commands
- Each command group has its own directory:
  - `auth/`: Authentication commands (login, logout, status)
    - `apikey/`: API key management (create, list, revoke)
  - `daemon/`: Daemon management
  - `server/`: Server lifecycle commands
  - `space/`: Space management operations
  - `shell/`: Interactive shell mode
  - `update/`: Self-update functionality
  - `version/`: Version information
- `root.go` registers all commands

### Core Logic (`/core/`)
- `client.go`: gRPC client singleton for server communication
- `auth.go`: Authentication logic with keyring integration
- `space.go`: Space management operations
- `stream.go`: Event streaming functionality with EventReceiver
- `keyring.go`: Secure credential storage (tokens and API keys)
- `apikey.go`: API key generation and management
- `config/constants.go`: Centralized configuration constants

### Daemon (`/daemon/`)
- `daemon.go`: Main daemon process that manages server lifecycle
- `taskmanager.go`: Schedules and manages background tasks
- `daemon_client.go`: Client for daemon communication

### Tasks (`/tasks/`)
- Background tasks executed by the daemon
- `autoapprove.go`: Auto-approves space join requests
- `server.go`: Manages server startup/shutdown

## Key Dependencies

- `github.com/anyproto/anytype-heart v0.41.2`: The middleware server
- `github.com/spf13/cobra v1.8.1`: CLI framework
- `google.golang.org/grpc v1.73.0`: gRPC communication
- `github.com/zalando/go-keyring`: Secure credential storage
- `github.com/cheggaaa/mb/v3 v3.2.0`: Message batching queue for event handling

## Important Notes

1. **Two-Process Architecture**: The CLI requires both a daemon process and the middleware server to be running
2. **Keyring Integration**: Authentication tokens are stored securely in the system keyring
3. **gRPC Communication**: All server interaction happens via gRPC on localhost:31007
4. **Event Streaming**: Uses server-sent events for real-time updates in auto-approval
5. **Version Management**: Version info is injected at build time via ldflags
6. **Self-Updating**: The CLI can update itself using the `anytype update` command
7. **API Keys**: Support for generating API keys for programmatic access

## Common Development Tasks

### Adding a New Command
1. Create a new directory under `/cmd/` for your command group
2. Create a file named after the command (e.g., `config.go` for config command) with a `NewXxxCmd()` function that returns `*cobra.Command`
3. Create subdirectories for each subcommand with their own files
4. Import subcommands with aliases matching the directory name
5. Register the command in `/cmd/root.go` using `NewXxxCmd()`
6. Implement core logic in `/core/` if needed

Directory structure follows the subcommand pattern:
```
cmd/
├── config/
│   ├── config.go      # Main command file (not cmd.go)
│   ├── get/
│   │   └── get.go     # Subcommand with NewGetCmd()
│   ├── set/
│   │   └── set.go     # Subcommand with NewSetCmd()
│   └── reset/
│       └── reset.go   # Subcommand with NewResetCmd()
```

Example main command file:
```go
// cmd/config/config.go
package config

import (
    "github.com/spf13/cobra"
    
    configGetCmd "github.com/anyproto/anytype-cli/cmd/config/get"
    configSetCmd "github.com/anyproto/anytype-cli/cmd/config/set"
    configResetCmd "github.com/anyproto/anytype-cli/cmd/config/reset"
)

func NewConfigCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "config <command>",
        Short: "Manage configuration",
    }
    
    cmd.AddCommand(configGetCmd.NewGetCmd())
    cmd.AddCommand(configSetCmd.NewSetCmd())
    cmd.AddCommand(configResetCmd.NewResetCmd())
    
    return cmd
}
```

### Working with the Daemon
- Daemon tasks go in `/tasks/`
- Register new tasks in `daemon/taskmanager.go`
- Use the `Task` interface for new task implementations

### Error Handling
- Client connection errors are handled in `core/client.go`
- Server startup errors are managed in `daemon/daemon.go:connectToServer`
- Use standard Go error wrapping with context

### API Key Management
- API keys are created and managed by the server via gRPC APIs
- The CLI provides commands to create, list, and revoke API keys
- Keys are generated server-side and can be used for programmatic access
