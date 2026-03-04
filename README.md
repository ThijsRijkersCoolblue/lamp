# lamp

A terminal emulator written in Go, built on top of [tcell](https://github.com/gdamore/tcell) and [pty](https://github.com/creack/pty).

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green?style=flat)

## Features

- Runs a real shell (`bash`) inside a PTY
- ANSI/VT100 escape sequence parsing
  - Cursor movement (absolute, relative, line-based)
  - SGR colors — 16 color, 256 color, and 24-bit truecolor
  - Erase in display and erase in line (`J`, `K`)
  - Insert/delete lines (`L`, `M`)
  - Scroll regions (`r`)
  - Bold, italic, dim, underline, blink, reverse
- UTF-8 character rendering
- Terminal resize support
- Works with interactive programs like `nvim`, `bash`

## Requirements

- Go 1.24+
- A Unix-like system (Linux, macOS)

## Installation

```bash
git clone https://github.com/thijsrijkers/lamp
cd lamp
go mod tidy
go build -o lamp .
./lamp
```

## How It Works

lamp spawns a shell process attached to a pseudo-terminal (PTY). Raw output from the PTY is parsed byte-by-byte in `ansi.ProcessOutput`, which interprets escape sequences and renders characters directly onto a `tcell.Screen`. Keyboard events are captured by tcell and forwarded back to the PTY as raw bytes.

## Keybindings

All standard terminal input is forwarded to the shell, including:

| Key | Action |
|-----|--------|
| Arrow keys | Cursor movement |
| `Ctrl+C` | Interrupt |
| `Ctrl+D` | EOF / logout |
| `Ctrl+Z` | Suspend |
| `Ctrl+L` | Clear screen |
| `F1`–`F12` | Function keys |
| `Home`, `End`, `PgUp`, `PgDn` | Navigation |

## Dependencies

- [tcell](https://github.com/gdamore/tcell) — terminal screen rendering
- [pty](https://github.com/creack/pty) — pseudo-terminal support

