# talswitcher

`talswitcher` is a command-line tool for managing and switching between different Talos contexts. It simplifies the process of selecting a Talos context from multiple talosconfig files and updates the active configuration.

## Features

- **Non-Interactive Mode**: Specify the desired context directly as an argument.
- **Interactive Mode**: Run the program with no arguments to get an interactive list of available contexts.
- **Automatic Configuration Update**: Copies the selected context's talosconfig to the appropriate location for Talos.
- **Duplicate Context Detection**: Ensures there are no duplicate context names across multiple talosconfig files.
- **Shell Completions**: Provides tab completion for available context names in compatible shells.

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

### Install via `go install` (Go 1.17+)

```bash
# Make sure $GOBIN is in your $PATH, or set it
go env -w GOBIN=~/.local/bin/ && export PATH="$PATH:$GOBIN"

# Install talswitcher
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
    mise run build
    ```

## Usage

`talswitcher` has two main subcommands:

### Context Subcommand

The `context` (or `ctx`) subcommand is used to switch between Talos contexts:

```bash
# Interactive mode
talswitcher context

# Switch to a specific context
talswitcher context my-context

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

|   Environment Variable   |      CLI Flag       |                     Description                     |                      Acceptable Values                      |   Default Value   |
| :----------------------: | :-----------------: | :-------------------------------------------------: | :---------------------------------------------------------: | :---------------: |
| `TALSWITCHER_LOG_LEVEL`  |    `--log-level`    |        Controls the logging verbosity level.        | `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` |      `info`       |
| `TALSWITCHER_LOG_FORMAT` |   `--log-format`    |           Controls the log output format.           |                       `json`, `text`                        |      `text`       |
|    `TALOSCONFIG_DIR`     | `--talosconfig-dir` |    Directory containing your talosconfig files.     |                  Any valid directory path                   |       `N/A`       |
|      `TALOSCONFIG`       |        `N/A`        | Path where the selected talosconfig will be copied. |                     Any valid file path                     | `~/.talos/config` |

When both the environment variable and CLI flag are set, the CLI flag takes precedence.

## Contributing

Contributions are welcome! Please fork the repository, make your changes, and submit a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
