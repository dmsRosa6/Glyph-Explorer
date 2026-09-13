# Scout

A small terminal file explorer built as a reference application for the **Glyph** TUI framework.

The goal is not to be a full-featured file manager. Scout is deliberately small, but it exercises a useful cross-section of Glyph:

- application lifecycle
- keyboard input and key bindings
- containers and layout
- borders and styling
- dynamic widget children
- selection and scrolling
- filesystem-driven state
- status/error feedback
- terminal-safe external file opening

## Run

From this directory:

```bash
go run . [path]
```

Examples:

```bash
go run .
go run . ~/Downloads
```

## Controls

| Key | Action |
|---|---|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` / `l` | Open directory/file |
| `←` / `h` | Parent directory |
| `Backspace` | Parent directory |
| `g` | First entry |
| `G` | Last entry |
| `.` | Toggle hidden files |
| `q` | Quit |

## Structure

```text
scout/
├── main.go          # Glyph application wiring + input bindings
├── explorer.go      # Explorer state and filesystem navigation
├── filesystem.go    # Filesystem primitives
├── view.go          # Glyph UI and row rendering
└── README.md
```

Scout intentionally keeps the domain state (`Explorer`) separate from the Glyph view. That makes it useful as a small integration/reference project while also making the framework API easier to evaluate.
# Glyph-Explorer
