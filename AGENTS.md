# AGENTS quick-ref

1. Build the binary:
   • `make build` (wraps `go build ./cmd/tmux-mcp`)

2. Run tests:
   • All packages   → `make test` or `go test ./... -v`
   • Single package → `go test ./internal/tmuxmcp`
   • Single test    → `go test -run '^TestPaneOutput$' ./internal/tmuxmcp`

3. Lint / format:
   • `go vet ./...`          – standard vetting
   • `golangci-lint run`     – if installed
   • `gofmt -s -w .` & `goimports -w .` – keep code tidy

4. Style guidelines:
   • Imports grouped: std / third-party / internal
   • Names: PascalCase for exported, camelCase for private, acronyms ALLCAPS (URI, ID)
   • Errors: return `error`; wrap with `%w` (`fmt.Errorf("open pane: %w", err)`) – never panic outside `main`
   • Keep funcs short (<50 lines); pass context as `ctx context.Context` when blocking/IO
   • Tests live in `*_test.go`, use `t.Helper()` inside helpers
   • Prefer interfaces over concrete types; avoid globals
   • Generics ok but favour clarity over cleverness

5. Cursor / Copilot rules:
   • None present (`.cursor/` or `.github/copilot-instructions.md` absent). Add as needed.
