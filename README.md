# talswitcher

`talswitcher` is a command-line tool for managing and switching between different Talos contexts. It simplifies the process of selecting a Talos context from multiple talosconfig files and updates the active configuration.

Just dump all your `talosconfig` into a single dir and let `talswitcher` manage them for you!

## Features

- **Multiple talosconfig files**: Manage multiple talosconfig files in a single directory without merging them
- **Context switching**: Switch between contexts from multiple config files
- **Interactive & non-interactive modes**: Select from a list or specify directly as an argument (with tab completion support!)

## Installation

### Download Precompiled Binaries

Precompiled binaries are available for various platforms. You can download the latest release from the [GitHub Releases page](https://github.com/mirceanton/talswitcher/releases/latest).

1. Download the appropriate binary for your system and extract the archive.
2. Make the extracted binary executable:

   ```bash
   chmod +x talswitcher
   ```

3. Move the binary to a directory in your PATH:

   ```bash
   mv talswitcher /usr/local/bin/
   ```

### Running via Docker

`talswitcher` is also available as a Docker container:

```bash
docker pull ghcr.io/mirceanton/talswitcher
```

### Install via homebrew

1. Add the tap

   ```bash
   brew tap mirceanton/taps
   ```

2. Install `talswitcher`

   ```bash
   brew install talswitcher
   ```

### Install via `go install`

```bash
go install github.com/mirceanton/talswitcher@main
```

### Build from Source

1. Clone the repository:

   ```bash
   git clone https://github.com/mirceanton/talswitcher
   cd talswitcher
   ```

2. Build the tool:

   ```bash
   go build -o talswitcher
   ```

## Usage

`talswitcher` has two main subcommands:

### Context Subcommand

The `context` (or `ctx`) subcommand is used to switch between Talos contexts:

```bash
# Interactive mode
talswitcher context

# Switch to a specific context
talswitcher ctx my-context

# Switch to previous context
talswitcher context -
```

### Shell Completions

The `completion` subcommand generates shell completion scripts:

```bash
# Generate completions for bash
talswitcher completion bash > /etc/bash_completion.d/talswitcher

# Generate completions for zsh
talswitcher completion zsh > ~/.zsh/completion/_talswitcher

# Generate completions for fish
talswitcher completion fish > ~/.config/fish/completions/talswitcher.fish

# Generate completions for powershell
talswitcher completion powershell > ~/talswitcher.ps1
```

## Configuration

You can configure `talswitcher` using environment variables or CLI flags. The following table outlines the available options:

| Option                | Flag                | Environment Variable | Default             | Description                                                       |
| --------------------- | ------------------- | -------------------- | ------------------- | ----------------------------------------------------------------- |
| Talosconfig Directory | `--talosconfig-dir` | `TALOSCONFIG_DIR`    | `~/.talos/configs/` | Directory containing your talosconfig files                       |
| Talosconfig           | `--talosconfig`     | `TALOSCONFIG`        | `~/.talos/config`   | Path to the currently active talosconfig file.                    |
| Log Level             | `--log-level`       | `LOG_LEVEL`          | `info`              | Logging verbosity (trace, debug, info, warn, error, fatal, panic) |
| Log Format            | `--log-format`      | `LOG_FORMAT`         | `text`              | Log output format (text, json)                                    |

When both the environment variable and CLI flag are set, the CLI flag takes precedence.

## Contributing

Contributions are welcome! Please fork the repository, make your changes, and submit a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
