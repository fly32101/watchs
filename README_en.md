# Watchs - File Change Monitor Tool

A file change monitoring tool based on DDD (Domain-Driven Design) architecture that can monitor file changes in a specified directory and execute commands when files change.

## Features

- Monitor file changes in specified directories (recursively)
- Filter by file type (supports multiple types)
- Exclude specific directories or files
- Execute commands when files change
- Support for configuration files and command-line parameters
- Support for generating configuration files via command line
- Interactive configuration wizard
- Clean, maintainable, and extensible code based on DDD architecture
- Extensible command-line interface using the Command Pattern
- Integrated GitHub Actions for automated builds and releases

## Project Architecture

The project adopts a DDD (Domain-Driven Design) architecture, divided into the following layers:

- **Domain Layer**: Contains core business logic and entities
  - `entity`: Domain entities such as configuration and file events
  - `service`: Domain service interfaces
  - `repository`: Repository interfaces

- **Application Layer**: Coordinates domain objects to complete user tasks
  - Application services, such as file monitoring service

- **Infrastructure Layer**: Provides technical implementations
  - `persistence`: Configuration persistence implementation
  - `watcher`: File monitoring and command execution implementation

- **Presentation Layer**: Handles user interaction
  - `cli`: Command-line interface using the Command Pattern

### Design Patterns

The project uses the following design patterns:

- **Command Pattern**: Encapsulates command-line operations as objects for extensibility and composability
- **Dependency Injection**: Injects dependencies through constructors to reduce coupling between components
- **Repository Pattern**: Abstracts data access logic, separating persistence from domain logic
- **Factory Method**: Creates complex objects, encapsulating object creation logic

## Automated Builds and Releases

The project uses GitHub Actions for automated builds and releases:

- **Continuous Integration (CI)**: Automatically tests and builds on Linux and Windows using Go 1.21 on every code push and PR
- **Automated Releases**: Automatically builds binaries and creates GitHub Releases when a new tag is created (e.g., v1.0.0)

### Releasing a New Version

To release a new version, simply create and push a new tag:

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

GitHub Actions will automatically build binaries and create a Release.

## Installation

### From GitHub Releases

Visit the [GitHub Releases](https://github.com/yourusername/watchs/releases) page and download the binary for your system.

### From Source

```bash
go install github.com/watchs/cmd/watchs@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/watchs.git
cd watchs
go build -o watchs ./cmd/watchs
```

## Usage

### View Help Information

```bash
watchs help
```

Or view help for a specific command:

```bash
watchs help <command_name>
watchs <command_name> --help
```

### View Version Information

```bash
watchs version
```

### Interactive Configuration

Use the interactive wizard to create a configuration file (recommended for new users):

```bash
watchs interactive
```

The wizard will guide you through all configuration options and allow you to start monitoring immediately.

### Generate Configuration File via Command Line

Use the `init` command to generate a configuration file:

```bash
watchs init -config watchs.json -dir ./ -types .go,.js -exclude vendor,node_modules -cmd "go run main.go"
```

Parameters:
- `-config`: Configuration file path (defaults to `watchs.json`)
- `-dir`: Directory to monitor (defaults to `./`)
- `-types`: File types to monitor, comma-separated
- `-exclude`: Paths to exclude, comma-separated
- `-cmd`: Command to execute when files change (defaults to `echo File updated`)
- `-force`: Whether to force overwrite existing configuration file

### Using a Configuration File

After creating a `watchs.json` configuration file, simply run:

```bash
watchs
```

Or specify a configuration file path:

```bash
watchs -config custom-watchs.json
```

You can also use the watch command (same as running directly):

```bash
watchs watch -config watchs.json
```

### Using Command Line Parameters

You can also run directly with command line parameters, without a configuration file:

```bash
watchs -dir ./ -types .go,.json -exclude vendor,node_modules,.git -cmd "go run main.go"
```

Or use the watch command:

```bash
watchs watch -dir ./ -types .go,.json -exclude vendor,node_modules,.git -cmd "go run main.go"
```

## Command Line Parameters

### Watch Command Parameters (watch)

- `-config`: Configuration file path (defaults to `watchs.json`)
- `-dir`: Directory to monitor (overrides configuration file)
- `-types`: File types to monitor, comma-separated (overrides configuration file)
- `-exclude`: Paths to exclude, comma-separated (overrides configuration file)
- `-cmd`: Command to execute when files change (overrides configuration file)
- `-debounce`: Debounce time in milliseconds (defaults to 500)

### Initialization Command Parameters (init)

- `-config`: Configuration file path (defaults to `watchs.json`)
- `-dir`: Directory to monitor (defaults to `./`)
- `-types`: File types to monitor, comma-separated
- `-exclude`: Paths to exclude, comma-separated
- `-cmd`: Command to execute when files change (defaults to `echo File updated`)
- `-force`: Whether to force overwrite existing configuration file

## Examples

### View Help

```bash
# Display all available commands
watchs help

# Display help for a specific command
watchs help interactive
watchs init --help
```

### Interactive Configuration

```bash
# Start the interactive configuration wizard
watchs interactive
```

### Generate Configuration File

```bash
# Generate default configuration file
watchs init

# Generate custom configuration file
watchs init -config frontend.json -dir ./frontend -types .js,.jsx,.ts,.tsx,.css -exclude node_modules -cmd "npm run build"
```

### Monitor File Changes

```bash
# Monitor using configuration file
watchs

# Monitor all .go files in the current directory, excluding vendor directory, and run tests when files change
watchs -dir ./ -types .go -exclude vendor -cmd "go test ./..."

# Monitor frontend project and automatically rebuild
watchs -dir ./frontend -types .js,.jsx,.ts,.tsx,.css -exclude node_modules -cmd "npm run build"
```

## Extending Commands

If you want to add a new command, simply implement the `Command` interface and register it during CLI initialization:

```go
// Implement the command interface
type MyCommand struct {
    // Dependencies
}

func (c *MyCommand) Name() string {
    return "mycommand"
}

func (c *MyCommand) Description() string {
    return "My custom command"
}

func (c *MyCommand) Execute(args []string) error {
    // Command implementation
    return nil
}

// Register during CLI initialization
registry.Register(NewMyCommand(...))
```

## Notes

- Commands are executed in the monitored directory
- If a command is a long-running process, it will be terminated and restarted when files change again
- Debounce mechanism is used to avoid frequent command execution 