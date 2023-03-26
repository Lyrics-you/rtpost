package rtp

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type RTP struct {
	// byte 0
	Version   uint8 // 2 bits
	Padding   uint8 // 1 bit
	Extension uint8 // 1 bit
	CSRC      uint8 // 4 bits, Contributing source

	// byte 1
	Marker       uint8 // 1 bit
	PlayloadType uint8 // 7 bit

	// byte 2-3
	SequenceNumber uint16 // 16 bits

	// byte 4-7
	Timestamp uint32 // 32 bits

	// byte 8-11
	SSRC uint32 // 32 bits, stream_id used here

	//CSRCItem
	CSRCItem []uint32
}

func (rtp *RTP) Decode(playload []byte) error {
	if len(playload) < 12 {
		return errors.New("can not decode to rtp protocol")
	}
	// byte 0
	// 2 bits
	rtp.Version = playload[0] & 0b11000000 >> 6
	// 1 bit
	rtp.Padding = playload[0] & 0b00100000 >> 5
	// 1 bit
	rtp.Extension = playload[0] & 0b00010000 >> 4
	// 4 bits
	rtp.CSRC = playload[0] & 0b00001111

	// byte 1
	// 1 bit
	rtp.Marker = playload[1] & 0b10000000 >> 7
	// 7 bits
	rtp.PlayloadType = playload[1] & 0b01111111

	// byte 2-3
	rtp.SequenceNumber = binary.BigEndian.Uint16(playload[2:4])
	// byte 4-7
	rtp.Timestamp = binary.BigEndian.Uint32(playload[4:8])
	// byte 8-11
	rtp.SSRC = binary.BigEndian.Uint32(playload[8:12])

	// CSRC Item
	for i := uint8(0); i < rtp.CSRC; i++ {
		if len(playload) >= int(16+4*i) {
			rtp.CSRCItem = append(rtp.CSRCItem, binary.BigEndian.Uint32(playload[12+4*i:16+4*i]))
		}
	}
	return nil
}

func (rtp *RTP) IsDynamicRTP() bool {
	if rtp.PlayloadType >= 96 && rtp.PlayloadType <= 172 {
		return true
	}
	return false
}

func (rtp *RTP) IsDynamicRTP116() bool {
	return rtp.PlayloadType == 116
}

func (rtp *RTP) IsUnassigned() bool {
	if rtp.PlayloadType >= 80 && rtp.PlayloadType <= 95 {
		return true
	}
	return false
}

func (rtp *RTP) IsTarget() bool {
	if rtp.IsDynamicRTP() || rtp.IsUnassigned() {
		return true
	}
	return false
}

func (rtp *RTP) ToCSRCItemString() string {
	if !rtp.IsTarget() {
		if rtp.CSRC == 0 {
			return "[]"
		} else {
			return "[...]"
		}
	}

	csrcItemStr := ""
	for _, i := range rtp.CSRCItem {
		if csrcItemStr != "" {
			csrcItemStr += fmt.Sprintf("%s, 0x%X", csrcItemStr, i)
		} else {
			csrcItemStr += fmt.Sprintf("0x%X", i)
		}
	}
	return fmt.Sprintf("[%s]", csrcItemStr)
}

func (rtp *RTP) ToSSRCString() string {
	return fmt.Sprintf("0x%X", rtp.SSRC)
}

func (rtp *RTP) ToWireSharkString() string {
	pt := fmt.Sprintf("PT=%v", rtp.PlayloadType)
	if rtp.IsUnassigned() {
		pt = fmt.Sprintf("PT=Unassigned-%v", rtp.PlayloadType)
	} else if rtp.IsDynamicRTP() {
		pt = fmt.Sprintf("PT=DynamicRTP-%v", rtp.PlayloadType)
	}
	ssrc := fmt.Sprintf("SSRC=0x%X", rtp.SSRC)
	seq := fmt.Sprintf("Seq=%v", rtp.SequenceNumber)
	time := fmt.Sprintf("Time=%v", rtp.Timestamp)
	csrc := fmt.Sprintf("CSRC=%v", rtp.CSRC)
	csrcItem := rtp.ToCSRCItemString()
	if rtp.Marker == 1 {
		return fmt.Sprintf("%s, %s, %s, %s, %s, %s, Mark", pt, ssrc, seq, time, csrc, csrcItem)
	}
	return fmt.Sprintf("%s, %s, %s, %s, %s, %s", pt, ssrc, seq, time, csrc, csrcItem)
}
