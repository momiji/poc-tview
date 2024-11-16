package ui

import (
	"fmt"
	"sync"
)

var appIsRunning bool
var appIsRunningLock sync.Mutex
var quitUI = make(chan any)
var StoppedUI = make(chan any)

func IfAppConsole(fn func(console bool)) {
	appIsRunningLock.Lock()
	defer appIsRunningLock.Unlock()
	fn(!appIsRunning)
}

func IfApp(fn func()) {
	IfAppConsole(func(console bool) {
		if !console {
			fn()
		}
	})
}

func IfConsole(fn func()) {
	IfAppConsole(func(console bool) {
		if console {
			fn()
		}
	})
}

func SwitchUI(console bool) {
	if console {
		appClose()
	} else {
		consoleClose()
	}
}

func RunUI(console bool) {
loop:
	for {
		select {
		case <-quitUI:
			break loop
		default:
		}
		if console {
			consoleRun()
		} else {
			suspendPrintUI()
			appInit()
			IfConsole(func() {
				appIsRunning = true
			})
			appRun()
			resumePrintUI()
		}
		console = !console
	}
	// first close App then Console, so we'll be in console mode at the end and normally resumePrintUI
	PrintUI("before appClose")
	appClose()
	PrintUI("before consoleClose")
	consoleClose()
	PrintUI("before appClosed.Wait()")
	// wait for app and console closed
	appClosed.Wait()
	PrintUI("before consoleClosed.Wait()")
	consoleClosed.Wait()
	// signal close
	close(StoppedUI)
}

func StopUI() {
	close(quitUI)
}

var suspendLock sync.RWMutex
var suspended bool

func suspendPrintUI() {
	suspendLock.Lock()
	defer suspendLock.Unlock()
	suspended = true
}

func resumePrintUI() {
	suspendLock.Lock()
	defer suspendLock.Unlock()
	suspended = false
}

func PrintUI(format string, a ...any) {
	suspendLock.RLock()
	defer suspendLock.RUnlock()
	if !suspended {
		fmt.Printf(format+"\n", a...)
	}
}
