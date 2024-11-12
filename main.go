// Demo code for the Table primitive.
package main

import (
	"example.com/m/ui"
	"math/rand"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var atomicReqId = atomic.Int32{}

func initTable() {
	s1 := "[aqua]pac-ret>ret[white] CONNECT [yellow]http://www.example.com[white] HTTP/1.1"
	s2 := "[aqua]socks[white] GET [yellow]http://www.google.com[white] HTTP/1.1"
	rows := []*ui.TrafficRow{
		ui.NewTraffic(0, s1),
		ui.NewTraffic(1, s2),
		ui.NewTraffic(2, s1),
		ui.NewTraffic(3, s2),
		ui.NewTraffic(4, s1),
		ui.NewTraffic(5, s2),
		ui.NewTraffic(6, s1),
		ui.NewTraffic(7, s2),
		ui.NewTraffic(8, s1),
		ui.NewTraffic(9, s2),
		ui.NewTraffic(10, s1),
		ui.NewTraffic(11, s1),
		ui.NewTraffic(12, s1),
		ui.NewTraffic(13, s2),
		ui.NewTraffic(14, s2),
	}
	for _, row := range rows {
		ui.Traffic.Add(row)
	}
	atomicReqId.Store(int32(len(rows)))
}

func updateTable() {
	// update table randomly by adding or removing rows and increasing or decreasing the values
	for i, row := range ui.Traffic.RowsCopy() {
		row.BytesSentPerSecond.IncrementBy(rand.Intn(10000000) * i * i)
		row.BytesReceivedPerSecond.IncrementBy(rand.Intn(100) * i * i)
	}
}

func updateForEver() {
	ticker := time.NewTicker(50 * time.Millisecond)
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
	syscall.SetNonblock(0, true)

	initTable()
	go updateForEver()

	// Start a dummy infinite logger
	go func() {
		for {
			ui.PrintUI("Log message: " + time.Now().String())
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Switch to console after 5 seconds
	go func() {
		time.Sleep(5 * time.Second)
		ui.PrintUI("Switching to console")
		ui.SwitchUI(true)
	}()

	// Switch to app after 2 seconds
	go func() {
		time.Sleep(500 * time.Millisecond)
		ui.SwitchUI(false)
	}()

	go ui.RunUI(true)
loop:
	for {
		select {
		case <-exitSignal:
			ui.StopUI()
		case <-ui.StoppedUI:
			break loop
		}
	}

	println("The END")

}
