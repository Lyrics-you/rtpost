package rtp

import (
	"fmt"
	"strings"
	"time"
)

type WireShark struct {
	ArrivalTime time.Time
	Delta       time.Duration
	SrcMac      string
	DstMac      string
	Source      string
	SrcPort     string
	Destination string
	DstPort     string
	Protocol    string
	Length      int
	Info        string
	RtpPacket   *RTP
}

func (ws *WireShark) ToString() string {
	arrival := ws.ArrivalTime.Format("2006-01-02 15:04:05.000000")
	return fmt.Sprintf("[%s] %s:%s->%s:%s %s[%d] {%s}", arrival, ws.Source, ws.SrcPort, ws.Destination, ws.DstPort, ws.Protocol, ws.Length, ws.Info)
}

func (ws *WireShark) ToHasDeltaString() string {
	arrival := ws.ArrivalTime.Format("2006-01-02 15:04:05.000000")
	return fmt.Sprintf("[%s][%.3fs] %s:%s->%s:%s %s[%d] {%s}", arrival, float64(ws.Delta.Milliseconds())/1000, ws.Source, ws.SrcPort, ws.Destination, ws.DstPort, ws.Protocol, ws.Length, ws.Info)
}

func (ws *WireShark) HasRtpSsrc(rtpSsrc string) bool {
	if ws.RtpPacket == nil {
		return false
	}
	return strings.Contains(ws.RtpPacket.ToSSRCString(), rtpSsrc)
}

func (ws *WireShark) HasRtpCsrc(rtpSsrc string) bool {
	if ws.RtpPacket == nil {
		return false
	}
	return strings.Contains(ws.RtpPacket.ToCSRCItemString(), rtpSsrc)
}

func (ws *WireShark) HasRtpPort(udpPorts []string) bool {
	for _, p := range udpPorts {
		if p == ws.DstPort || p == ws.SrcPort {
			return true
		}
	}
	return false
}

func (ws *WireShark) DstInPorts(udpPorts []string) bool {
	for _, p := range udpPorts {
		if p == ws.DstPort {
			return true
		}
	}
	return false
}

func (ws *WireShark) SrcInPorts(udpPorts []string) bool {
	for _, p := range udpPorts {
		if p == ws.SrcPort {
			return true
		}
	}
	return false
}
