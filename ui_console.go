package main

import (
	"golang.org/x/term"
	"io"
	"os"
)

var ConsoleClosed = NewManualResetEvent(true)
var closeConsole = NewManualResetEvent(false)
var consoleReader chan byte

func ConsoleRun() {
	ConsoleClosed.Reset()
	defer ConsoleClosed.Signal()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	if consoleReader == nil {
		consoleReader = ioReader(os.Stdin)
	}
	loop := true
	for loop {
		select {
		case <-closeConsole.Channel():
			closeConsole.Reset()
			loop = false
		case b, ok := <-consoleReader:
			if !ok {
				loop = false
				break
			}
			IfConsole(func() {
				switch b {
				case 'q', 'Q', '\x03':
					quit.Signal()
					loop = false
				case ' ', '\x1b', '\x09':
					loop = false
				}
			})
		}
	}
}

func ConsoleClose() {
	closeConsole.Signal()
}

func ioReader(reader io.Reader) chan byte {
	var err error
	b := make([]byte, 1)
	ch := make(chan byte)
	go func() {
		for {
			_, err = reader.Read(b)
			if err != nil {
				close(ch)
				return
			}
			ch <- b[0]
		}
	}()
	return ch
}
