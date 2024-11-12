// Demo code for the Table primitive.
package main

import (
	"math/rand"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var atomicReqId = atomic.Int32{}

func initTable() {
	trafficTable.Table = []*Traffic{
		NewTraffic(0, "[pac-ret>ret] CONNECT http://www.example.com HTTP 1.1"),
		NewTraffic(1, "http://www.google.com"),
		NewTraffic(2, "http://www.yahoo.com"),
		NewTraffic(3, "http://www.example.com"),
		NewTraffic(4, "http://www.google.com"),
		NewTraffic(5, "http://www.yahoo.com"),
		NewTraffic(6, "http://www.example.com"),
		NewTraffic(7, "http://www.google.com"),
		NewTraffic(8, "http://www.yahoo.com"),
		NewTraffic(9, "http://www.example.com"),
		NewTraffic(10, "http://www.google.com"),
		NewTraffic(11, "http://www.yahoo.com"),
		NewTraffic(12, "http://www.example.com"),
		NewTraffic(13, "http://www.google.com"),
		NewTraffic(14, "http://www.yahoo.com"),
	}
	atomicReqId.Store(3)
}

func updateTable() {
	// update table randomly by adding or removing rows and increasing or decreasing the values
	for i, row := range trafficTable.Table {
		row.BytesSentPerSecond.IncrementBy(rand.Intn(10000000) * i * i)
		row.BytesReceivedPerSecond.IncrementBy(rand.Intn(100) * i * i)
	}
}

func updateForEver() {
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			updateTable()
		}
	}
}

var exitSignal = make(chan os.Signal, 1)

func main() {
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	initTable()
	go updateForEver()

	// Start a dummy infinite logger
	go func() {
		for {
			Log("Log message: " + time.Now().String())
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Switch to console every 5 seconds
	go func() {
		for {
			time.Sleep(5 * time.Second)
			Log("Switching to console")
			AppClose()
		}
	}()

	// Switch to app after 10 seconds
	go func() {
		time.Sleep(10 * time.Second)
		ConsoleClose()
	}()

	runApp := false
loop:
	for {
		select {
		case <-exitSignal:
			break loop
		case <-quit.Channel():
			break loop
		default:
		}
		if runApp {
			SuspendLogs()
			AppInit()
			IfConsole(func() {
				appIsRunning = true
			})
			AppRun()
			ResumeLogs()
		} else {
			ConsoleRun()
		}
		runApp = !runApp
	}
	AppClose()
	ConsoleClose()
	AppClosed.Wait()
	ConsoleClosed.Wait()
}
