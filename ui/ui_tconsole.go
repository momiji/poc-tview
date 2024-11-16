package ui

import (
	"example.com/m/term"
	"os"
	"time"
)

var consoleClosed = NewManualResetEvent(true)
var closeConsole = NewManualResetEvent(false)
var consoleInited = false
var consoleChan = make(chan byte)

func consoleRun() {
	if !consoleInited {
		go readConsole(consoleChan)
		consoleInited = true
	}
	consoleClosed.Reset()
	defer consoleClosed.Signal()
	closeConsole.Reset()
	// backup term state and restore it on return
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	// clean consoleChan
	for len(consoleChan) > 0 {
		<-consoleChan
	}
	// console loop
	b := make([]byte, 1)
	ticker := time.NewTicker(100 * time.Millisecond)
	ticker.Stop()
	loop := true
	for loop {
		select {
		case <-closeConsole.c:
			closeConsole.Reset()
			loop = false
		case bb := <-consoleChan:
			IfConsole(func() {
				switch bb {
				case 'q', 'Q', '\x03':
					close(quitUI)
					loop = false
				case ' ', '\x1b', '\x09':
					loop = false
				case '\x0a', '\x0d':
					PrintUI("")
				}
			})
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

func consoleClose() {
	closeConsole.Signal()
}

func readConsole(c chan byte) {
	b := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(b)
		if err != nil {
			continue
		}
		if n > 0 {
			IfAppConsole(func(console bool) {
				if console {
					c <- b[0]
				} else {
					appKeyNoLock(b[0])
				}
			})

		}
	}
}
