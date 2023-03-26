package rtp

import (
	"fmt"
)

var (
	packageLength = int32(0)
	packageLoss   = int32(0)
)

type WirePackets []*WireShark

// func isDuplicate(ws1, ws2 *WireShark) bool {
// 	return ws1.Source == ws2.Source &&
// 		ws1.SrcPort == ws2.SrcPort &&
// 		ws1.Destination == ws2.Destination &&
// 		ws1.DstPort == ws2.DstPort &&
// 		ws1.RtpPacket.SequenceNumber == ws2.RtpPacket.SequenceNumber
// }

// func (wp *WirePackets) Deduplicate() *WirePackets {
// 	var deduplicated WirePackets
// 	for _, pkt := range *wp {
// 		unique := true
// 		for _, dpkt := range deduplicated {
// 			if isDuplicate(dpkt, pkt) {
// 				unique = false
// 				break
// 			}
// 		}
// 		if unique {
// 			deduplicated = append(deduplicated, pkt)
// 		}
// 	}
// 	wp = &deduplicated
// 	return &deduplicated
// }

func (wp *WirePackets) Deduplicate() *WirePackets {
	deduplicated := WirePackets{}
	if len(*wp) == 0 {
		return &deduplicated
	}
	last := int32((*wp)[0].RtpPacket.SequenceNumber) - 1
	for _, p := range *wp {
		if CheckSequenceNumber(last, int32(p.RtpPacket.SequenceNumber)) {
			deduplicated = append(deduplicated, p)
			last = int32(p.RtpPacket.SequenceNumber)
		} else {
			// repeat
			continue
		}
	}
	return &deduplicated
}

func (wp *WirePackets) Delta() *WirePackets {
	if len(*wp) == 0 {
		return wp
	}

	for i, p := range *wp {
		if i != 0 {
			p.Delta = p.ArrivalTime.Sub((*wp)[i-1].ArrivalTime)
		}
	}
	return wp
}

func (wp *WirePackets) RemoveNotDynamicRTP() *WirePackets {
	rewp := WirePackets{}
	for _, p := range *wp {
		if p.RtpPacket.IsDynamicRTP() {
			rewp = append(rewp, p)
		}
	}
	wp = &rewp
	return &rewp
}

func (wp *WirePackets) RemoveNotUnassigned() *WirePackets {
	rewp := WirePackets{}
	for _, p := range *wp {
		if p.RtpPacket.IsDynamicRTP() {
			rewp = append(rewp, p)
		}
	}
	wp = &rewp
	return &rewp
}

func (wp *WirePackets) RemoveNotTargetType() *WirePackets {
	rewp := WirePackets{}
	for _, p := range *wp {
		if p.RtpPacket.IsTarget() {
			rewp = append(rewp, p)
		}
	}
	wp = &rewp
	return &rewp
}

func (wp *WirePackets) DeduplicateGroup() *[]*WirePackets {
	m := make(map[string]WirePackets)
	deduplicated_group := []*WirePackets{}
	if len(*wp) == 0 {
		return &deduplicated_group
	}
	for _, pkt := range *wp {
		key := fmt.Sprintf("%s:%s->%s:%s-%v", pkt.Source, pkt.SrcPort, pkt.Destination, pkt.DstPort, pkt.RtpPacket.SSRC)
		if _, ok := m[key]; !ok {
			m[key] = WirePackets{}
		} else {
			m[key] = append(m[key], pkt)
		}
	}
	for g := range m {
		mg := m[g]
		deduplicated_group = append(deduplicated_group, mg.Deduplicate())
	}
	return &deduplicated_group
}

func (wp *WirePackets) DoGroup() *[]*WirePackets {
	m := make(map[string]WirePackets)
	group := []*WirePackets{}
	if len(*wp) == 0 {
		return &group
	}
	for _, pkt := range *wp {
		key := fmt.Sprintf("%s:%s->%s:%s-%v", pkt.Source, pkt.SrcPort, pkt.Destination, pkt.DstPort, pkt.RtpPacket.SSRC)
		if _, ok := m[key]; !ok {
			m[key] = WirePackets{}
		} else {
			m[key] = append(m[key], pkt)
		}
	}
	for g := range m {
		mg := m[g]
		group = append(group, &mg)
	}
	return &group
}

func CheckSequenceNumber(s1, s2 int32) bool {
	return (s2 > s1 && s2-s1 < 32768) || (s1 > s2 && s1-s2 > 32768)
}

func DiffSequenceNumber(s1, s2 int32) int32 {
	diff := int32(0)
	if s2 >= s1 {
		diff = s2 - s1
	} else {
		diff = (65536 + s2) - s1
	}
	return diff
}

func (wp *WirePackets) ExpectedPacketsLength() (length int32) {
	if len(*wp) == 0 {
		return 0
	}
	last := int32((*wp)[0].RtpPacket.SequenceNumber) - 1
	length = int32(0)

	for _, p := range *wp {
		if CheckSequenceNumber(last, int32(p.RtpPacket.SequenceNumber)) {
			length += DiffSequenceNumber(last, int32(p.RtpPacket.SequenceNumber))
			last = int32(p.RtpPacket.SequenceNumber)
		} else {
			// repeat
			continue
		}
	}
	return
}

func (wp *WirePackets) ActualLoss() (loss int32) {
	if len(*wp) == 0 {
		return 0
	}
	last := int32((*wp)[0].RtpPacket.SequenceNumber) - 1
	loss = int32(0)
	for _, p := range *wp {
		if CheckSequenceNumber(last, int32(p.RtpPacket.SequenceNumber)) {
			// lost
			if DiffSequenceNumber(last, int32(p.RtpPacket.SequenceNumber)) > 1 {
				// fmt.Println(p.RtpPacket.ToSSRCString(), last, "-", p.RtpPacket.SequenceNumber)
				loss += DiffSequenceNumber(last, int32(p.RtpPacket.SequenceNumber)) - 1
			}
			last = int32(p.RtpPacket.SequenceNumber)
		} else {
			// repeat
			continue
		}
	}
	return
}

func (wp *WirePackets) ExpectedPacketsLengthAndLoss() (length int32, loss int32) {
	if len(*wp) == 0 {
		return 0, 0
	}
	last := int32((*wp)[0].RtpPacket.SequenceNumber) - 1
	length = int32(0)
	loss = int32(0)
	for _, p := range *wp {
		if CheckSequenceNumber(last, int32(p.RtpPacket.SequenceNumber)) {
			length += DiffSequenceNumber(last, int32(p.RtpPacket.SequenceNumber))
			// lost
			if DiffSequenceNumber(last, int32(p.RtpPacket.SequenceNumber)) > 1 {
				// fmt.Println(p.RtpPacket.ToSSRCString(), last, "-", p.RtpPacket.SequenceNumber)
				loss += DiffSequenceNumber(last, int32(p.RtpPacket.SequenceNumber)) - 1
			}
			last = int32(p.RtpPacket.SequenceNumber)
		} else {
			// repeat
			continue
		}
	}
	return
}

func (wp *WirePackets) LostSequenceNumbers() (int32, *[]string) {
	if len(*wp) == 0 {
		return 0, &[]string{}
	}
	lostString := []string{}
	last := int32((*wp)[0].RtpPacket.SequenceNumber) - 1
	line := "----------------------------"
	lasfInfo := ""
	loss := int32(0)
	for _, p := range *wp {
		if CheckSequenceNumber(last, int32(p.RtpPacket.SequenceNumber)) {
			// lost
			if DiffSequenceNumber(last, int32(p.RtpPacket.SequenceNumber)) > 1 {
				loss += DiffSequenceNumber(last, int32(p.RtpPacket.SequenceNumber)) - 1
				// fmt.Println(p.RtpPacket.ToSSRCString(), last, "-", p.RtpPacket.SequenceNumber)
				lostString = append(lostString, line)
				lostString = append(lostString, lasfInfo)
				start := last + 1
				for DiffSequenceNumber(start, int32(p.RtpPacket.SequenceNumber)) > 0 {
					lostString = append(lostString, fmt.Sprintf("|- lost Seq=%d", start))
					start += 1
				}
				lostString = append(lostString, p.ToString())
			}
			last = int32(p.RtpPacket.SequenceNumber)
			lasfInfo = p.ToString()
		} else {
			// repeat
			continue
		}
	}
	lostString = append(lostString, line)
	return loss, &lostString
}

func (wp *WirePackets) LossRate() float64 {
	rate := float64(int(packageLength)-len(*wp)) / float64(packageLength)
	return rate
}

func (wp *WirePackets) Loss() int32 {
	_, delta := wp.ExpectedPacketsLengthAndLoss()
	return delta
}

func (wp *WirePackets) LossRateInfo() string {
	if len(*wp) == 0 {
		return "packet is zero"
	}
	wp0 := (*wp)[0]
	packageLength, packageLoss = wp.ExpectedPacketsLengthAndLoss()
	return fmt.Sprintf("%s:%s->%s:%s {SSRC:%s CSRC:%s} (recv:%d, expect:%d, loss:%d, rate:%.2f%%)", wp0.Source, wp0.SrcPort, wp0.Destination, wp0.DstPort, wp0.RtpPacket.ToSSRCString(), wp0.RtpPacket.ToCSRCItemString(), len(*wp), packageLength, packageLoss, wp.LossRate()*100)
}
