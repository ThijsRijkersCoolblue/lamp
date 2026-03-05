# lamp

A small terminal emulator written in Go, built on top of [tcell](https://github.com/gdamore/tcell), [pty](https://github.com/creack/pty), and [Fyne](https://fyne.io).

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green?style=flat)
![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux-lightgrey?style=flat)

## How It Works

Lamp spawns a shell process attached to a pseudo-terminal (PTY). Raw output from the PTY is parsed byte-by-byte in `ansi.ProcessOutput`, which interprets escape sequences and writes characters into a `tcell.SimulationScreen`, an in-memory cell buffer. Keyboard events from Fyne are mapped to tcell key events and forwarded back to the PTY as raw bytes.

## Installation

### Run in terminal
```bash
git clone https://github.com/thijsrijkers/lamp
cd lamp
go mod tidy
go build -o lamp .
./lamp
```

### Install as a native macOS app
```bash
go mod tidy
make install-macos
```

This builds a `Lamp.app` bundle, ad-hoc codesigns it, and installs it to `/Applications`. You can then launch it from Finder, Spotlight, or pin it to your Dock.

### Install on Linux
```bash
make install-linux
```
*(NOTE: This is not tested properly yet)*

Installs the binary to `/usr/local/bin` and adds a `.desktop` entry to your application launcher.

## Dependencies

- [tcell](https://github.com/gdamore/tcell) — terminal screen buffer and ANSI parsing
- [pty](https://github.com/creack/pty) — pseudo-terminal support
- [fyne](https://fyne.io) — native windowing and GPU-accelerated rendering
- [golang.org/x/image](https://pkg.go.dev/golang.org/x/image) — font rendering with OpenType support

## Requirements

- Go 1.24+
- macOS 10.13+ or Linux
- Xcode Command Line Tools (macOS): `xcode-select --install`

