package main

import (
	"golang.org/x/term"
	"os"
	"time"
)

var ConsoleClosed = NewManualResetEvent(true)
var closeConsole = NewManualResetEvent(false)

func ConsoleRun() {
	ConsoleClosed.Reset()
	defer ConsoleClosed.Signal()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	b := make([]byte, 1)
	ticker := time.NewTicker(100 * time.Millisecond)
	loop := true
	for loop {
		select {
		case <-closeConsole.Channel():
			closeConsole.Reset()
			loop = false
		case <-ticker.C:
			n, err := os.Stdin.Read(b)
			if err != nil {
				continue
			}
			if n > 0 {
				IfConsole(func() {
					switch b[0] {
					case 'q', 'Q', '\x03':
						close(quitUI)
						loop = false
					case ' ', '\x1b', '\x09':
						loop = false
					}
				})
			}
		}
	}
}

func ConsoleClose() {
	closeConsole.Signal()
}
