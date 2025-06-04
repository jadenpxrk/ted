# Ted

Ted is the fastest way to get an answer in the terminal.

## Features

- **Agent Mode**: Generate a single command from natural language with confirmation
- **Ask Mode**: Get multiple command suggestions for your questions
- **History**: Interactive browser for your command history
- **Settings**: Easy configuration management for API keys and preferences

## Installation

### From Source

```bash
git clone https://github.com/jadenpxrk/ted
cd ted
go build -o ted
sudo mv ted /usr/local/bin/
```

## Setup

Before using Ted, you'll need to configure your Gemini API key:

1. Get your API key from [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Run the settings command:

```bash
ted settings
```

3. Enter your API key when prompted

## Usage

### Agent Mode

Generate and execute a single command:

```bash
ted agent how to make a python virtual environment
# Output: You can create a Python virtual environment by running 'python3 -m venv myenv'
# [y/N] to execute
```

### Ask Mode

Get multiple command suggestions:

```bash
ted ask how to find large files
# Output:
# 1. 'find . -type f -size +100M' - Find files larger than 100MB in current directory
# 2. 'du -h --max-depth=1 | sort -hr' - Show directory sizes sorted by size
# 3. 'ls -lah | sort -k5 -hr' - List files sorted by size
# Select an option (1-3) or press Enter to exit:
```

### History

Browse your command history:

```bash
ted history
```

Navigate through your recent commands (last 5 entries), view details, and manage your history.

### Settings

Configure your preferences:

```bash
ted settings
```

## Config

Ted stores its configuration in `~/.ted/`:

- `config.yaml` - API keys and model settings
- `history.db` - Command history (BoltDB database, limited to last 5 entries)

## Available Models

- `gemini-2.0-flash` (default)
- `gemini-2.0-flash-lite`
- `gemini-2.5-pro-preview-05-06`
- `gemini-2.5-flash-preview-05-20`

## Examples

```bash
# System administration
ted agent check disk usage
ted ask how to monitor system resources

# File operations
ted agent compress a folder to zip
ted ask how to find files modified today

# Development
ted agent create a git repository
ted ask how to check git status

# Network
ted agent check open ports
ted ask how to test network connectivity
```

## License

MIT License - see LICENSE file for details.

## Project Structure

```
ted/
├── cmd/                   # Cobra CLI commands
│   ├── agent.go           # Agent command (single command generation)
│   ├── ask.go             # Ask command (multiple suggestions)
│   ├── history.go         # History browsing
│   ├── settings.go        # Configuration management
│   └── root.go            # Root command and help
├── internal/
│   ├── colors/            # Centralized color and styling
│   │   └── colors.go      # All UI colors and styles
│   ├── config/            # Configuration management
│   │   └── config.go      # Viper-based config handling
│   ├── gemini/            # Google Gemini AI integration
│   │   └── gemini.go      # API client and response parsing
│   ├── history/           # Command history management
│   │   └── history.go     # History storage and retrieval
│   └── ui/                # User interface components
│       └── ui.go          # Bubble Tea confirmation dialogs
├── main.go                # Application entry point
└── go.mod                 # Go module definition
```
