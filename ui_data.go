package main

import (
	"github.com/enterprizesoftware/rate-counter"
	"time"
)

type Traffic struct {
	ReqId                  int32
	Url                    string
	BytesSent              int64
	BytesReceived          int64
	BytesSentPerSecond     *ratecounter.Rate
	BytesReceivedPerSecond *ratecounter.Rate
}

type TrafficTable struct {
	Table []*Traffic
}

func NewTraffic(reqId int32, url string) *Traffic {
	return &Traffic{
		ReqId:                  reqId,
		Url:                    url,
		BytesSentPerSecond:     ratecounter.New(100*time.Millisecond, 5*time.Second),
		BytesReceivedPerSecond: ratecounter.New(100*time.Millisecond, 5*time.Second),
	}
}

var trafficTable = TrafficTable{}
