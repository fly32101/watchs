# Watchs - File Change Monitor

A file change monitoring tool based on DDD (Domain-Driven Design) architecture, which can monitor file changes in specified directories and execute specified commands when files change.

## Features

* Monitor file changes in specified directories (recursive)
* Filter by file type (supports multiple types)
* Exclude specific directories or files
* Execute specified commands when files change
* Support configuration files and command line parameters
* Generate configuration files via command line
* Interactive configuration wizard
* Based on DDD architecture, clear code structure, easy to maintain and extend
* Implement extensible command-line interface using Command Pattern
* Integrated GitHub Actions for automated building and releasing

## Latest Version

The current latest version is [v1.0.0](https://github.com/fly32101/watchs/releases/tag/v1.0.0), automatically built and released through GitHub Actions.

## Project Architecture

The project adopts DDD (Domain-Driven Design) architecture, divided into the following layers:

* **Domain Layer**: Contains core business logic and entities
  * `entity`: Domain entities, such as configuration and file events
  * `service`: Domain service interfaces
  * `repository`: Repository interfaces
* **Application Layer**: Coordinates domain objects to complete user tasks
  * Application services, such as file monitoring service
* **Infrastructure Layer**: Provides technical implementation
  * `persistence`: Configuration persistence implementation
  * `watcher`: File monitoring and command execution implementation
* **Presentation Layer**: Handles user interaction
  * `cli`: Command-line interface, implemented using the Command Pattern

### Design Patterns

The project uses the following design patterns:

* **Command Pattern**: Encapsulates command-line operations as objects, enabling extensibility and composability
* **Dependency Injection**: Injects dependencies through constructors, reducing coupling between components
* **Repository Pattern**: Abstracts data access logic, separating persistence from domain logic
* **Factory Method**: Creates complex objects, encapsulating object creation logic

## Automated Build and Release

The project uses GitHub Actions for automated building and releasing:

* **Continuous Integration (CI)**: Automatically tests and builds on Linux and Windows using Go 1.21 for every code push and PR
* **Automatic Release**: Automatically builds binaries and creates GitHub Release when a new tag (such as v1.0.0) is created

### Release a New Version

To release a new version, simply create and push a new tag:

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

GitHub Actions will automatically build binaries and create a Release.

## Installation

### Install from GitHub Releases

Visit the [GitHub Releases page](https://github.com/fly32101/watchs/releases) and download the binary suitable for your system:
- Linux: `watchs_Linux_x86_64.tar.gz` (Intel/AMD) or `watchs_Linux_arm64.tar.gz` (ARM)
- Windows: `watchs_Windows_x86_64.zip`
- macOS: `watchs_Darwin_x86_64.tar.gz` (Intel) or `watchs_Darwin_arm64.tar.gz` (Apple Silicon)

### Install from Source

```bash
go install github.com/fly32101/watchs/cmd/watchs@latest
```

Or compile from source:

```bash
git clone https://github.com/fly32101/watchs.git
cd watchs
go build -o watchs ./cmd/watchs
```

## Usage

### View Help Information

```bash
watchs help
```

Or view help information for a specific command:

```bash
watchs help <command-name>
watchs <command-name> --help
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

The wizard will guide you through all configuration options and can optionally start monitoring immediately.

### Generate Configuration File via Command Line

Use the `init` command to generate a configuration file:

```bash
watchs init -config watchs.json -dir ./ -types .go,.js -exclude vendor,node_modules -cmd "go run main.go"
```

Parameter description:

* `-config`: Configuration file path (default is `watchs.json`)
* `-dir`: Directory to monitor (default is `./`)
* `-types`: File types to monitor, comma-separated
* `-exclude`: Paths to exclude, comma-separated
* `-cmd`: Command to execute when files change (default is `echo File updated`)
* `-force`: Whether to forcibly overwrite existing configuration file

### Use Configuration File

After creating the `watchs.json` configuration file, run directly:

```bash
watchs
```

Or specify the configuration file path:

```bash
watchs -config custom-watchs.json
```

You can also use the watch command (same as running directly):

```bash
watchs watch -config watchs.json
```

### Use Command Line Parameters

You can also run directly through command line parameters without a configuration file:

```bash
watchs -dir ./ -types .go,.json -exclude vendor,node_modules,.git -cmd "go run main.go"
```

Or use the watch command:

```bash
watchs watch -dir ./ -types .go,.json -exclude vendor,node_modules,.git -cmd "go run main.go"
```

## Command Line Parameters

### Watch Command Parameters (watch)

* `-config`: Configuration file path (default is `watchs.json`)
* `-dir`: Directory to monitor (overrides configuration file)
* `-types`: File types to monitor, comma-separated (overrides configuration file)
* `-exclude`: Paths to exclude, comma-separated (overrides configuration file)
* `-cmd`: Command to execute when files change (overrides configuration file)
* `-debounce`: Debounce time in milliseconds (default is 500)

### Initialization Command Parameters (init)

* `-config`: Configuration file path (default is `watchs.json`)
* `-dir`: Directory to monitor (default is `./`)
* `-types`: File types to monitor, comma-separated
* `-exclude`: Paths to exclude, comma-separated
* `-cmd`: Command to execute when files change (default is `echo File updated`)
* `-force`: Whether to forcibly overwrite existing configuration file

## Examples

### View Help

```bash
# Display all available commands
watchs help

# Display help information for a specific command
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

# Monitor all .go files in the current directory, excluding the vendor directory, and run tests when files change
watchs -dir ./ -types .go -exclude vendor -cmd "go test ./..."

# Monitor frontend project and automatically rebuild
watchs -dir ./frontend -types .js,.jsx,.ts,.tsx,.css -exclude node_modules -cmd "npm run build"
```

## Extend Commands

If you want to add a new command, simply implement the `Command` interface and register it during `CLI` initialization:

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

* Commands are executed in the monitored directory
* If the command is a long-running process, it will be terminated and restarted when files change again
* Uses debounce mechanism to avoid frequent command execution

## License

This project is licensed under the [MIT License](LICENSE).

## Contribution

Issues and Pull Requests are welcome.

## Links

- [GitHub Repository](https://github.com/fly32101/watchs)
- [Releases Page](https://github.com/fly32101/watchs/releases) 