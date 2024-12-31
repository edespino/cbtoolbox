
# Apache Cloudberry (Incubating) Toolbox

The Apache Cloudberry Toolbox (`cbtoolbox`) is a command-line utility that provides various tools and utilities for managing and monitoring Apache Cloudberry installations.

## Overview

This toolbox provides a collection of utilities to:
- Gather system and database environment information
- Analyze core dump files for diagnostics
- Support database administration tasks
- Monitor database performance
- Facilitate cluster management

## Installation

### Prerequisites
- Go 1.19 or later
- Access to an Apache Cloudberry installation
- Proper GPHOME environment setup

### Building from Source
You can build the project either manually using `go build` or using the included Makefile for a streamlined process.

#### Manual Build
```bash
git clone https://github.com/edespino/cbtoolbox.git
cd cbtoolbox
go build
```

#### Using Makefile
```bash
git clone https://github.com/edespino/cbtoolbox.git
cd cbtoolbox
make build
```

### Installing the Binary
```bash
go install github.com/edespino/cbtoolbox@latest
```

## Usage

Basic command structure:
```bash
cbtoolbox [command] [flags]
```

### Available Commands

1. **sysinfo**
   - Displays system and database environment information
   - See [sysinfo documentation](./cmd/sysinfo/README.md) for details

2. **coreinfo**
   - Analyzes core dump files for diagnostic purposes
   - Executes basic or detailed GDB commands based on debug symbol availability
   - See [coreinfo documentation](./cmd/coreinfo/README.md) for details

### Global Flags
- `--help, -h`: Display help information

## Environment Requirements

The toolbox requires specific environment setup:

### GPHOME
- Points to the Apache Cloudberry installation directory
- Required for database-specific functionality
- Example: `/usr/local/cloudberry-db-1.6.0`

## Makefile Usage

The Makefile simplifies common tasks like building, testing, and cleaning the project. Below are the available targets:

- **`make build`**: Builds the project binary and places it in the `build/` directory.
- **`make test`**: Runs all tests in verbose mode.
- **`make test-cover`**: Runs all tests with coverage reporting.
- **`make lint`**: Runs code linting using `golangci-lint`. Ensure it is installed beforehand.
- **`make clean`**: Cleans up build artifacts.
- **`make run`**: Runs the application binary.

Example usage:
```bash
# Build the project
make build

# Run tests
make test

# Clean up build artifacts
make clean

# Run the built application
make run
```

## Command Documentation

Detailed documentation for each command:
- [Command Package](./cmd/README.md)
  - [Sysinfo Command](./cmd/sysinfo/README.md)
  - [Coreinfo Command](./cmd/coreinfo/README.md)

## Project Structure

```
cbtoolbox/
├── cmd/                  # Command implementations
│   ├── root.go           # Root command
│   ├── sysinfo/          # Sysinfo command
│   └── coreinfo/         # Coreinfo command
├── main.go               # Application entry point
├── Makefile              # Build and task automation
└── README.md             # Project documentation
```

## Development

### Testing

Run all tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-cover
```

### Adding New Commands

1. Create a new package under `cmd/`
2. Implement the command using cobra.Command
3. Register your command in cmd/root.go
4. Add comprehensive tests for your command
5. Document your command in its package README.md

### Code Style

Follow standard Go conventions:
- Run `go fmt` before committing
- Use `golint` for style checking
- Follow Go Code Review Comments

## Error Handling

The toolbox implements consistent error handling:
- Command-specific error codes
- Descriptive error messages
- Appropriate exit codes

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

### Guidelines
- Follow existing code structure
- Include tests for new functionality
- Update documentation
- Add command documentation
- Follow Go best practices

## License

Licensed under the Apache License, Version 2.0. See LICENSE for details.

## Support

For support and issues:
- GitHub Issues: [edespino/cbtoolbox/issues](https://github.com/edespino/cbtoolbox/issues)
