package main

import (
	"fmt"
	"sync"
)

var suspendLock sync.RWMutex
var suspended bool

func SuspendLogs() {
	suspendLock.Lock()
	defer suspendLock.Unlock()
	suspended = true
}

func ResumeLogs() {
	suspendLock.Lock()
	defer suspendLock.Unlock()
	suspended = false
}

func Log(s string) {
	suspendLock.RLock()
	defer suspendLock.RUnlock()
	if !suspended {
		fmt.Printf("%s\r\n", s)
	}
}
