package main

import (
	"github.com/dustin/go-humanize"
	"github.com/enterprizesoftware/rate-counter"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
	"strings"
	"time"
)

var screen tcell.Screen
var app *tview.Application
var table *tview.Table
var AppClosed = NewManualResetEvent(true)

func BytesFormat(rate *ratecounter.Rate) string {
	return humanize.Comma(int64(rate.Total()))
}

func RateFormat(rate *ratecounter.Rate) string {
	return strings.ReplaceAll(humanize.IBytes(uint64(rate.RatePer(1*time.Second))), "i", "")
}

func setCell(i, j int, s string, w int, left bool, newRow bool) {
	align := tview.AlignRight
	if left {
		align = tview.AlignLeft
	}
	if w > 0 {
		if len(s) < w {
			if left {
				s += strings.Repeat(" ", w-len(s))
			} else {
				s = strings.Repeat(" ", w-len(s)) + s
			}
		}
	} else if w < 0 {
		if len(s) > -w {
			s = s[:-w-1] + "…"
		} else if len(s) < -w {
			if left {
				s += strings.Repeat(" ", -w-len(s))
			} else {
				s = strings.Repeat(" ", -w-len(s)) + s
			}
		}
	}
	s = " " + s + " "
	if newRow {
		color := tcell.ColorWhite
		if i == 0 {
			color = tcell.ColorYellow
		}
		table.SetCell(i, j, tview.NewTableCell(s).SetAlign(align).SetTextColor(color))
	} else {
		table.GetCell(i, j).Text = s
	}
}

func setRow(row int, new bool, urlWidth int, reqId string, url string, bytesSent string, bytesReceived string, bytesSentPerSecond string, bytesReceivedPerSecond string) {
	setCell(row, 0, reqId, 5, false, new)
	setCell(row, 1, url, -urlWidth, true, new)
	setCell(row, 2, bytesReceived, 15, false, new)
	setCell(row, 3, bytesSent, 15, false, new)
	setCell(row, 4, bytesReceivedPerSecond, 7, false, new)
	setCell(row, 5, bytesSentPerSecond, 7, false, new)
}

func AppInit() {
	if app != nil {
		return
	}

	var err error

	// Create screen
	if screen, err = tcell.NewScreen(); err != nil {
		panic(err)
	}

	// Create table
	table = tview.NewTable()
	table.
		ScrollToBeginning().
		SetBorders(false).
		SetFixed(1, 0).
		SetSelectable(false, false).
		SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)).
		SetSeparator(tview.Borders.Vertical)
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return nil
	})

	setRow(0, true, 0, "ID", "URL", "RECV", "SENT", "RECV/S", "SENT/S")

	// Create application
	app = tview.NewApplication().EnableMouse(false)
	app.SetRoot(table, true).SetScreen(screen)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				IfApp(func() {
					appIsRunning = false
				})
				quit.Signal()
				app.Stop()
			case ' ':
				IfApp(func() {
					appIsRunning = false
				})
				app.Stop()
			}
		case tcell.KeyCtrlC:
			IfApp(func() {
				appIsRunning = false
			})
			quit.Signal()
			app.Stop()
		case tcell.KeyEsc, tcell.KeyTab:
			IfApp(func() {
				appIsRunning = false
			})
			app.Stop()
		default:
		}
		return nil
	})

	go AppUpdate()
}

func AppRun() {
	AppClosed.Reset() // Run application
	defer AppClosed.Signal()
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func AppClose() {
	IfApp(func() {
		app.QueueEvent(tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone))
	})
}

func AppUpdate() {
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-quit.Channel():
			return
		case <-ticker.C:
			IfApp(func() {
				app.QueueUpdateDraw(func() {
					// update the table with the new data
					screenWidth, screenHeight := screen.Size()
					totalWidth := 5 + 15 + 15 + 7 + 7 + 15 + 2
					urlWidth := screenWidth - totalWidth
					if urlWidth < 20 {
						urlWidth = 20
					}
					if table.GetRowCount() == 0 {
						setRow(0, true, 0, "ID", "URL", "RECV", "SENT", "RECV/S", "SENT/S")
					}
					for i, row := range trafficTable.Table {
						if i+1 >= screenHeight {
							break
						}
						newRow := i+1 >= table.GetRowCount()
						setRow(i+1, newRow, urlWidth,
							strconv.Itoa(int(row.ReqId)),
							row.Url,
							BytesFormat(row.BytesSentPerSecond),
							BytesFormat(row.BytesReceivedPerSecond),
							RateFormat(row.BytesSentPerSecond),
							RateFormat(row.BytesReceivedPerSecond))
					}
					// remove any extra rows
					for i := table.GetRowCount() - 1; i > len(trafficTable.Table); i-- {
						table.RemoveRow(i)
					}
					// remove hidden rows
					for i := screenHeight; i < table.GetRowCount(); i++ {
						table.RemoveRow(i)
					}
				})
			})
		}
	}
}
