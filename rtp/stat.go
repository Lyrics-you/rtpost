package rtp

import "fmt"

type WirePackets []*WireShark

func (wp WirePackets) LossRate() float64 {
	l := len(wp)
	total := wp[l-1].SequenceNumber - wp[0].SequenceNumber + 1
	rate := float64(wp.Loss()) / float64(total)
	return rate
}

func (wp WirePackets) Loss() int {
	l := len(wp)
	total := wp[l-1].SequenceNumber - wp[0].SequenceNumber + 1
	delta := int(total) - l
	return delta
}

func (wp WirePackets) LossRateInfo() string {
	if len(wp) == 0 {
		return "packet is zero"
	}
	return fmt.Sprintf("%s:%s->%s:%s (recv:%d, loss:%d, rate:%.2f%%)", wp[0].Source, wp[0].SrcPort, wp[0].Destination, wp[0].DstPort, len(wp), wp.Loss(), wp.LossRate()*100)
}
