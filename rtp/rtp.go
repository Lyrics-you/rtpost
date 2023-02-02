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
	return nil
}

func (rtp *RTP) ToWireSharkString() string {
	pt := "Unassigned"
	if rtp.PlayloadType == 116 {
		pt = "PT=DynamicRTP-116"
	}
	ssrc := fmt.Sprintf("SSRC=0x%X", rtp.SSRC)
	seq := fmt.Sprintf("Seq=%v", rtp.SequenceNumber)
	time := fmt.Sprintf("Time=%v", rtp.Timestamp)
	if rtp.Marker == 1 {
		return fmt.Sprintf("%s, %s, %s, %s, Mark", pt, ssrc, seq, time)
	}
	return fmt.Sprintf("%s, %s, %s, %s", pt, ssrc, seq, time)
}
