# DeadCodeFinder

DeadCodeFinder is a lightweight, offline CLI tool that scans a repository for dead code across Python, JavaScript/TypeScript, Go, and Rust. It reports unused functions and can optionally remove them.

## Installation

```bash
go install github.com/scobieworks/deadcodefinder@latest
```

## Usage

```bash
deadcodefinder .          # Scan current directory
deadcodefinder . --json   # Output JSON
deadcodefinder . --fix    # Remove unused definitions (dry-run by default)
deadcodefinder . --dry-run --fix  # Show what would be removed
```

## Flags

- `-exclude` : Comma-separated list of directories to exclude.
- `-dry-run` : Show what would be removed without modifying files.
- `-fix`     : Automatically remove unused definitions.
- `-json`    : Output report in JSON format.

## License

MIT

<!-- ORION-MONETIZATION:START -->
## Support

Download the CLI and consider a small donation to support continued maintenance: https://paypal.me/Damonwill.
<!-- ORION-MONETIZATION:END -->
