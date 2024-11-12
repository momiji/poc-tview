package main

import "sync"

var appIsRunning bool
var appIsRunningLock sync.Mutex
var quit = NewManualResetEvent(false)

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

func switchApp(console bool) {
	if console {
		AppClose()
	} else {
		ConsoleClose()
	}
}
