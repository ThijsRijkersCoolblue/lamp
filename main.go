package main

import (
	"lamp/ansi"
	"lamp/window"
	"log"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gdamore/tcell/v2"
)

func main() {
	// Start shell in PTY
	cmd := exec.Command("bash")
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Fatal(err)
	}
	defer ptmx.Close()

	screen := window.CreateScreen()

	width, height := screen.Size()
	pty.Setsize(ptmx, &pty.Winsize{
		Cols: uint16(width),
		Rows: uint16(height),
	})

	cursorX, cursorY := 0, 0
	ansiState := ansi.NewState(height)

	// Read PTY output
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				return
			}
			ansi.ProcessOutput(screen, buf[:n], &cursorX, &cursorY, ansiState)
		}
	}()

	// Write helper
	writeToPTY := func(b []byte) {
		ptmx.Write(b)
	}

	// Event loop
	for {
		event := screen.PollEvent()
		switch ev := event.(type) {
		case *tcell.EventResize:
			screen.Sync()
			newWidth, newHeight := screen.Size()
			ansiState.ScrollBottom = newHeight - 1
			pty.Setsize(ptmx, &pty.Winsize{
				Cols: uint16(newWidth),
				Rows: uint16(newHeight),
			})
		default:
			window.HandleEvent(screen, ev, writeToPTY)
		}
	}
}
