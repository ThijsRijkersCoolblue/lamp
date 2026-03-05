## Features

- Runs a real shell (`bash`) inside a PTY
- Native macOS window via Fyne — no Terminal.app required
- ANSI/VT100 escape sequence parsing
  - Cursor movement (absolute, relative, line-based)
  - SGR colors — 16 color, 256 color, and 24-bit truecolor
  - Erase in display and erase in line (`J`, `K`)
  - Insert/delete lines (`L`, `M`)
  - Scroll regions (`r`)
  - Bold, italic, dim, underline, blink, reverse
- UTF-8 character rendering
- Retina display support (2× pixel density)
- Terminal resize support
- Works with interactive programs like `nvim`, `bash`

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