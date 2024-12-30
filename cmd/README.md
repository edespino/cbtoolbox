# Command Package

The `cmd` package implements the core command-line interface for the Apache Cloudberry Toolbox. It serves as the central point for registering and managing all subcommands available in the toolbox.

## Overview

This package provides:
- Root command initialization and configuration
- Subcommand registration and management
- Common command-line flags and utilities
- Error handling and output formatting

## Package Structure

```
cmd/
├── root.go           # Root command implementation
├── root_test.go      # Root command tests
└── sysinfo/         # Sysinfo subcommand package
```

## Root Command

The root command (`rootCmd`) serves as the entry point for the Apache Cloudberry Toolbox. It:
- Provides the main help and usage information
- Manages subcommand registration
- Handles global flags and configuration

### Usage

```bash
cbtoolbox [command] [flags]
```

### Available Commands

1. sysinfo
   - Displays system and database environment information
   - See [sysinfo documentation](./sysinfo/README.md) for details

2. help
   - Displays help about any command
   - Usage: `cbtoolbox help [command]`

3. completion
   - Generates shell completion scripts
   - Usage: `cbtoolbox completion [bash|zsh|fish|powershell]`

### Global Flags

- `--help, -h`: Display help information about any command

## Implementation Details

### Command Registration

New commands are registered in the root command's initialization:

```go
func init() {
    rootCmd.AddCommand(sysinfo.Cmd)
}
```

### Error Handling

The package implements consistent error handling across all commands:
- Command-specific errors are propagated to the root command
- Errors are formatted consistently for user output
- Exit codes indicate success (0) or failure (non-zero)

## Development

### Adding New Commands

1. Create a new package under `cmd/` for your command
2. Implement your command using the cobra.Command structure
3. Register your command in `cmd/root.go`
4. Add comprehensive tests for your command
5. Document your command in its package README.md

### Testing

Run the test suite:
```bash
go test -v ./...
```

Run with test coverage:
```bash
go test -v -cover ./...
```

The test suite includes:
- Unit tests for command registration
- Integration tests for command execution
- Error handling validation
- Help text verification

## Error Codes

Commands should return appropriate error codes:
- 0: Success
- 1: General error (command failure, invalid arguments)
- Other codes may be defined for specific command failures

## Contributing

When contributing new commands:
1. Follow the existing command structure
2. Include comprehensive tests
3. Provide detailed documentation
4. Ensure proper error handling
5. Maintain consistent output formatting

## License

Licensed under the Apache License, Version 2.0. See LICENSE for details.
