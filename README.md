# Apache Cloudberry (Incubating) Toolbox

The Apache Cloudberry Toolbox (`cbtoolbox`) is a command-line utility that provides various tools and utilities for managing and monitoring Apache Cloudberry installations.

## Overview

This toolbox provides a collection of utilities to:
- Gather system and database environment information
- Support database administration tasks
- Monitor database performance
- Facilitate cluster management

## Installation

### Prerequisites
- Go 1.19 or later
- Access to a Apache Cloudberry installation
- Proper GPHOME environment setup

### Building from Source
```bash
git clone https://github.com/edespino/cbtoolbox.git
cd cbtoolbox
go build
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

1. sysinfo
   - Displays system and database environment information
   - See [sysinfo documentation](./cmd/sysinfo/README.md) for details

### Global Flags
- `--help, -h`: Display help information

## Environment Requirements

The toolbox requires specific environment setup:

### GPHOME
- Points to the Apache Cloudberry installation directory
- Required for database-specific functionality
- Example: `/usr/local/cloudberry-db-1.6.0`

## Command Documentation

Detailed documentation for each command:
- [Command Package](./cmd/README.md)
  - [Sysinfo Command](./cmd/sysinfo/README.md)

## Project Structure

```
cbtoolbox/
├── cmd/                  # Command implementations
│   ├── root.go          # Root command
│   └── sysinfo/         # Sysinfo command
├── main.go              # Application entry point
└── README.md           # Project documentation
```

## Development

### Testing

Run all tests:
```bash
go test -v ./...
```

Run tests with coverage:
```bash
go test -v -cover ./...
```

### Adding New Commands

1. Create a new package under `cmd/`
2. Implement the command using cobra.Command
3. Register in cmd/root.go
4. Add tests and documentation
5. Update command documentation

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

