# tmux-mcp

**tmux-mcp** is an [MCP](https://modelcontextprotocol.dev/) (Model Context Protocol) server that exposes rich, structured context from your running tmux server so agents and tools can reason about what’s happening in your terminal.

## What it does

| Capability | Description |
| ----------- | ----------- |
| List panes/sessions/windows | JSON resources describing every pane, window and session across all tmux servers |
| Capture pane output        | Tool that returns the last *N* lines from a specific pane |
| Tag panes (`@mcptag`)      | Adds semantic labels to panes; unique-tag helper ensures only one pane has a given tag |
| Find pane by tag           | Fuzzy search tool that returns a single pane when its tag is (uniquely) matched |

These resources are available to any MCP-aware LLM agent, enabling workflows like:

* "Show me logs from the backend pane"
* "Switch to the pane tagged `db` and run migrations"
* Automatically surface failing test output to the assistant when a test pane fails

## Quick start

```bash
# build binary → bin/tmux-mcp
make build

# run the server (defaults are fine for local use)
./bin/tmux-mcp
```

The server registers its tools & resources with MCP on start-up; compatible agents will auto-discover them.

## Usage examples

### Capture pane output

```jsonc
{
  "tool": "get_pane_output",
  "args": {
    "session_name": "mysesh",
    "window_index": 0,
    "pane_id": "%5",
    "lines": 100
  }
}
```

### Tagging sessions for the MCP
In order to tag a pane so that is it is identified by the MCP, you must add a tag to the pane using the `@mcptag` prefix. This can be done manually in tmux by using the following configuration:

`.tmux.conf`
```tmux
# Prefix + T → Tag a pane with a unique identifier
bind-key T command-prompt -p "Unique tag for pane %: " "run 'tmux-mcptag.sh %1 #{pane_id}'"

# Prefix + C-T → clear @mcptag on current pane
bind-key C-T set -p -u @mcptag
```


`tmux-mcptag.sh`
```bash
#!/usr/bin/env bash

# Script to tag a tmux pane with a unique identifier
TAG="$1"
PANE_ID="$2"

# Clear tag from any pane that has the same tag
tmux list-panes -a -F "#{pane_id} #{@mcptag}" | while read id t; do
	if [ "$t" = "$TAG" ]; then
		tmux set-option -p -u -t "$id" @mcptag
	fi
done

# Set tag on current pane
tmux set-option -p -t "$PANE_ID" @mcptag "$TAG"
tmux display-message "Pane tagged as: $TAG"

```

### Find pane by tag

```jsonc
{
  "tool": "find_pane_by_tag",
  "args": { "tag": "api" }
}
```

## Development

* Go 1.22+
* Make targets: `make build`, `make test`, `make tidy`
* Lint: `go vet ./...` or `golangci-lint run`
* Tests use a temporary tmux server; ensure `tmux` is installed.

```bash
make test   # run all tests
```

### Project layout

```
cmd/tmux-mcp       – main entry-point
internal/tmuxmcp   – MCP resources & tools
internal/exec      – thin wrapper around gotmux for exec tests
internal/desc      – long-form descriptions bundled with resources
```

## Contributing

PRs & issues welcome! Open an issue to discuss larger ideas first.

## License

MIT © 2025 Nick Apisani
