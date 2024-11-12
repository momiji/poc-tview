package main

import (
	"fmt"
	"sync"
)

var appIsRunning bool
var appIsRunningLock sync.Mutex
var quitUI = make(chan any)
var StoppedUI = make(chan any)

func ifAppConsole(running bool, fn func()) {
	appIsRunningLock.Lock()
	defer appIsRunningLock.Unlock()
	if appIsRunning == running {
		fn()
	}
}

func IfApp(fn func()) {
	ifAppConsole(true, fn)
}

func IfConsole(fn func()) {
	ifAppConsole(false, fn)
}

func SwitchUI(console bool) {
	if console {
		AppClose()
	} else {
		ConsoleClose()
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
			ConsoleRun()
		} else {
			SuspendPrintUI()
			AppInit()
			IfConsole(func() {
				appIsRunning = true
			})
			AppRun()
			ResumePrintUI()
		}
		console = !console
	}
	// first close App then Console, so we'll be in console mode at the end and normally ResumePrintUI
	AppClose()
	ConsoleClose()
	// wait for app and console closed
	AppClosed.Wait()
	ConsoleClosed.Wait()
	// signal close
	close(StoppedUI)
}

func StopUI() {
	close(quitUI)
}

var suspendLock sync.RWMutex
var suspended bool

func SuspendPrintUI() {
	suspendLock.Lock()
	defer suspendLock.Unlock()
	suspended = true
}

func ResumePrintUI() {
	suspendLock.Lock()
	defer suspendLock.Unlock()
	suspended = false
}

func PrintUI(format string, a ...any) {
	suspendLock.RLock()
	defer suspendLock.RUnlock()
	if !suspended {
		fmt.Printf(format+"\r\n", a...)
	}
}
