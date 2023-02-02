package rtp

import (
	"fmt"
	"strings"
)

type WireShark struct {
	SrcMac         string
	DstMac         string
	Source         string
	SrcPort        string
	Destination    string
	DstPort        string
	Protocol       string
	Length         int
	Info           string
	SequenceNumber uint16
}

func (ws *WireShark) ToString() string {
	return fmt.Sprintf("%s:%s->%s:%s %s[%d] {%s}", ws.Source, ws.SrcPort, ws.Destination, ws.DstPort, ws.Protocol, ws.Length, ws.Info)
}

func (ws *WireShark) HasRtpSsrc(rtpSsrc string) bool {
	return strings.Contains(ws.Info, rtpSsrc)
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
