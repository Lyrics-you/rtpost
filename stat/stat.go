package stat

import (
	"fmt"
	"rtpost/rtp"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	udpPorts []string = []string{"8089", "8090", "9600", "9601", "9602", "9603", "9604", "9605", "18089", "18090", "19600", "19601", "19602", "19603", "19604", "19605"}
	handle   *pcap.Handle
	err      error

	sllLayer layers.LinuxSLL
	ethLayer layers.Ethernet
	ipLayer  layers.IPv4
	udpLayer layers.UDP
)

func DoStat(pcapFile, rtpSsrc string, visual bool) {
	// Open file instead of device
	handle, err = pcap.OpenOffline(pcapFile)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Loop through packets in file
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	var srcPacket []*rtp.WireShark
	var dstPacket []*rtp.WireShark
	// start := time.Now()
	for packet := range packetSource.Packets() {
		// wireshark := ParsePackets(packet)
		// faster
		wireshark := DecodeLayers(packet)

		if wireshark.HasRtpSsrc(rtpSsrc) {
			if visual {
				fmt.Println(wireshark.ToString())
			}
			if wireshark.SrcInPorts(udpPorts) {
				srcPacket = append(srcPacket, wireshark)
			} else if wireshark.DstInPorts(udpPorts) {
				dstPacket = append(dstPacket, wireshark)
			}
			// break
		}

		// printPacketInfo(packet)
	}
	// cost := time.Since(start)
	// fmt.Println(cost)

	d := (rtp.WirePackets)(dstPacket)
	if len(d) != 0 {
		fmt.Println(d.LossRateInfo())
	}

	s := (rtp.WirePackets)(srcPacket)
	if len(s) != 0 {
		fmt.Println(s.LossRateInfo())
	}
}
