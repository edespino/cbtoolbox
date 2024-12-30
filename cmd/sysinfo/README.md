# Sysinfo Command

The `sysinfo` command is a component of the Apache Cloudberry Toolbox that gathers and displays detailed system and database environment information. It provides a comprehensive overview of both the host system and the Apache Cloudberry installation.

## Overview

The command collects information about:

### System Information
- Operating System and version
- System architecture
- Hostname
- Kernel version
- CPU count
- Memory statistics (Total, Free, Available, Cached, Buffers)

### Database Information (when GPHOME is set)
- GPHOME path validation
- PostgreSQL build configuration
- PostgreSQL server version
- Apache Cloudberry version

## Prerequisites

- Linux-based operating system
- GPHOME environment variable set to Apache Cloudberry installation directory
- Access to `/proc/meminfo` for memory statistics
- Execution permissions for `pg_config` and `postgres` binaries

## Usage

```bash
cbtoolbox sysinfo [flags]
```

### Flags
- `--format`: Output format (yaml or json). Default: "yaml"
- `--help`: Display help information

### Examples

1. Default output (YAML format):
```bash
cbtoolbox sysinfo
```

2. JSON format output:
```bash
cbtoolbox sysinfo --format=json
```

## Output Format

### YAML Output Example
```yaml
os: linux
architecture: amd64
hostname: cdw
kernel: Linux 4.18.0-553.el8_10.x86_64
os_version: Rocky Linux 8.10 (Green Obsidian)
cpus: 16
memory_stats:
  Buffers: 5.1 MiB
  Cached: 982.1 MiB
  MemAvailable: 60.5 GiB
  MemFree: 60.1 GiB
  MemTotal: 61.6 GiB
GPHOME: /usr/local/cloudberry-db-1.6.0
pg_config_configure:
  - --prefix=/usr/local/cloudberry-db
  - --disable-external-fts
  - --enable-gpcloud
postgres_version: postgres (Cloudberry Database) 14.4
gp_version: postgres (Cloudberry Database) 1.6.0 build 1
```

### JSON Output Example
```json
{
  "os": "linux",
  "architecture": "amd64",
  "hostname": "cdw",
  "kernel": "Linux 4.18.0-553.el8_10.x86_64",
  "os_version": "Rocky Linux 8.10 (Green Obsidian)",
  "cpus": 16,
  "memory_stats": {
    "Buffers": "5.1 MiB",
    "Cached": "982.1 MiB",
    "MemAvailable": "60.5 GiB",
    "MemFree": "60.1 GiB",
    "MemTotal": "61.6 GiB"
  },
  "GPHOME": "/usr/local/cloudberry-db-1.6.0",
  "pg_config_configure": [
    "--prefix=/usr/local/cloudberry-db",
    "--disable-external-fts",
    "--enable-gpcloud"
  ],
  "postgres_version": "postgres (Cloudberry Database) 14.4",
  "gp_version": "postgres (Cloudberry Database) 1.6.0 build 1"
}
```

## Error Handling

The command handles various error conditions:

1. GPHOME not set:
   - Displays available system information
   - Returns error about missing GPHOME
   - Exits with non-zero status

2. Missing executables:
   - Reports specific missing components
   - Continues collecting available information

3. Invalid format:
   - Returns error message
   - Shows valid format options

## Implementation Details

### Features
- Concurrent collection of system information for improved performance
- Thread-safe data gathering and error handling
- Automatic unit conversion for memory statistics (KiB, MiB, GiB)
- Graceful degradation when components are unavailable

### Memory Statistics
- Memory values are automatically converted to human-readable format
- Units are adjusted based on size (KiB, MiB, GiB)
- Original values from /proc/meminfo are preserved during conversion

## Development

### Testing
Run the test suite:
```bash
go test -v ./...
```

Run with test coverage:
```bash
go test -v -cover ./...
```

### Test Coverage
The test suite includes:
- Unit tests for all major functions
- Integration tests for system information gathering
- Concurrent execution testing
- Error condition validation
- Mock environment testing

## License

Licensed under the Apache License, Version 2.0. See LICENSE for details.
