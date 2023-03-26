package cycle

import (
	"fmt"
	"rtpost/logger"
	"rtpost/rtp"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// decode_as_entry: udp.port,18089,(none),RExecute
// decode_as_entry: udp.port,19000,(none),RTP
// decode_as_entry: udp.port,19600,(none),RTP
// decode_as_entry: udp.port,19601,(none),RTP
// decode_as_entry: udp.port,19602,(none),RTP
// decode_as_entry: udp.port,19603,(none),RTP
// decode_as_entry: udp.port,19604,(none),RTP
// decode_as_entry: udp.port,19605,(none),RTP
// decode_as_entry: udp.port,8089,(none),RTP
// decode_as_entry: udp.port,9000,(none),RTP
// decode_as_entry: udp.port,9600,OMRON FINS,RTP
// decode_as_entry: udp.port,9601,(none),RTP
// decode_as_entry: udp.port,9602,(none),RTP
// decode_as_entry: udp.port,9603,(none),RTP
// decode_as_entry: udp.port,9604,(none),RTP
// decode_as_entry: udp.port,9605,(none),RTP

// unix: yum install libpcap-devel

var (
	udpPorts []string = []string{"8089", "8090", "9600", "9601", "9602", "9603", "9604", "9605", "18089", "18090", "19600", "19601", "19602", "19603", "19604", "19605"}
	handle   *pcap.Handle
	err      error
	log      = logger.Logger()

	sllLayer layers.LinuxSLL
	ethLayer layers.Ethernet
	ipLayer  layers.IPv4
	udpLayer layers.UDP
)

func PrintPacketInfo(packet *gopacket.Packet) {
	// Let's see if the packet is an ethernet packet
	ethernetLayer := (*packet).Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		fmt.Println("Ethernet layer detected.")
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		fmt.Println("Source MAC: ", ethernetPacket.SrcMAC)
		fmt.Println("Destination MAC: ", ethernetPacket.DstMAC)
		// Ethernet type is typically IPv4 but could be ARP or other
		fmt.Println("Ethernet type: ", ethernetPacket.EthernetType)
		fmt.Println()
	}
	// Let's see if the packet is IP (even though the ether type told us)
	ipLayer := (*packet).Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		fmt.Println("IPv4 layer detected.")
		ip, _ := ipLayer.(*layers.IPv4)
		// IP layer variables:
		// Version (Either 4 or 6)
		// IHL (IP Header Length in 32-bit words)
		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
		// Checksum, SrcIP, DstIP
		fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
		fmt.Println("Protocol: ", ip.Protocol)
		fmt.Println()
	}
	// Let's see if the packet is TCP
	tcpLayer := (*packet).Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		fmt.Println("TCP layer detected.")
		tcp, _ := tcpLayer.(*layers.TCP)
		// TCP layer variables:
		// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
		// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
		fmt.Printf("From port %d to %d\n", tcp.SrcPort, tcp.DstPort)
		fmt.Println("Sequence number: ", tcp.Seq)
		fmt.Println()
	}
	// Let's see if the packet is UDP
	udpLayer := (*packet).Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		fmt.Println("UDP layer detected.")
		udp, _ := udpLayer.(*layers.UDP)
		// UDP layer variables:
		fmt.Printf("From port %d to %d\n", udp.SrcPort, udp.DstPort)
		fmt.Println()
	}
	// Iterate over all layers, printing out each layer type
	fmt.Println("All packet layers:")
	for _, layer := range (*packet).Layers() {
		fmt.Println("- ", layer.LayerType())
	}
	// When iterating through packet.Layers() above,
	// if it lists Payload layer then that is the same as
	// this applicationLayer. applicationLayer contains the payload
	applicationLayer := (*packet).ApplicationLayer()
	if applicationLayer != nil {
		fmt.Println("Application layer/Payload found.")
		fmt.Printf("%s\n", applicationLayer.Payload())
		// Search for a string inside the payload
		if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
			fmt.Println("HTTP found!")
		}
	}
	// Check for errors
	if err := (*packet).ErrorLayer(); err != nil {
		fmt.Println("Error decoding some part of the packet:", err)
	}
}

func DecodeLayers(packet *gopacket.Packet) *rtp.WireShark {
	wireshark := rtp.WireShark{}

	// get time from metadata
	meta := (*packet).Metadata()
	wireshark.ArrivalTime = meta.Timestamp
	// fmt.Println(meta.Timestamp)

	// if the packet is an ethernet packet
	ethernetLayer := (*packet).Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		// fmt.Println("Source MAC: ", ethernetPacket.SrcMAC)
		// fmt.Println("Destination MAC: ", ethernetPacket.DstMAC)
		// Ethernet type is typically IPv4 but could be ARP or other
		// fmt.Println("Ethernet type: ", ethernetPacket.EthernetType)
		wireshark.SrcMac = string(ethernetPacket.SrcMAC)
		wireshark.DstMac = string(ethernetPacket.DstMAC)
		wireshark.Length = len(ethernetPacket.Payload)
	}

	// if the packet is an linux sll packet
	linuxSllLayer := (*packet).Layer(layers.LayerTypeLinuxSLL)
	if linuxSllLayer != nil {
		linuxSllPacket, _ := linuxSllLayer.(*layers.LinuxSLL)
		wireshark.Length = len(linuxSllPacket.Payload)
		wireshark.DstMac = linuxSllPacket.Addr.String()
	}

	// ipLayer := (*packet).Layer(layers.LayerTypeIPv6)
	// if ipLayer != nil {
	// 	ip, _ := ipLayer.(*layers.IPv6)
	// 	fmt.Println(ip.SrcIP.String())
	// 	fmt.Println(ip.DstIP.String())
	// }

	// if the packet is IP (even though the ether type told us)
	ipLayer := (*packet).Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		wireshark.Source = ip.SrcIP.String()
		wireshark.Destination = ip.DstIP.String()
	}

	// if the packet is UDP
	udpLayer := (*packet).Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		wireshark.SrcPort = udp.SrcPort.String()
		wireshark.DstPort = udp.DstPort.String()
	}

	// if it lists Payload layer then that is the same as
	// this applicationLayer. applicationLayer contains the payload
	applicationLayer := (*packet).ApplicationLayer()
	if applicationLayer != nil {
		playload := applicationLayer.Payload()
		rtp := rtp.RTP{}
		err := rtp.Decode(playload)
		if err != nil {
			wireshark.Protocol = "UDP"
			wireshark.Info = fmt.Sprintf("%s -> %s Len(%v)", wireshark.SrcPort, wireshark.DstPort, len(playload))
		}
		wireshark.Protocol = "RTP"
		// wireshark.SequenceNumber = rtp.SequenceNumber
		// wireshark.Timestamp = rtp.Timestamp
		wireshark.RtpPacket = &rtp
		wireshark.Info = rtp.ToWireSharkString()
	}
	return &wireshark
}

func ParsePackets(packet *gopacket.Packet) *rtp.WireShark {
	wireshark := rtp.WireShark{}

	// get time from metadata
	meta := (*packet).Metadata()
	wireshark.ArrivalTime = meta.Timestamp

	var parser *gopacket.DecodingLayerParser
	if (*packet).Layers()[0].LayerType() == layers.LayerTypeLinuxSLL {
		parser = gopacket.NewDecodingLayerParser(
			layers.LayerTypeLinuxSLL,
			&sllLayer,
			&ipLayer,
			&udpLayer,
		)
	} else if (*packet).Layers()[0].LayerType() == layers.LayerTypeEthernet {
		parser = gopacket.NewDecodingLayerParser(
			layers.LayerTypeEthernet,
			&ethLayer,
			&ipLayer,
			&udpLayer,
		)
	}

	foundLayerTypes := []gopacket.LayerType{}
	err := parser.DecodeLayers((*packet).Data(), &foundLayerTypes)
	for _, layerType := range foundLayerTypes {
		if layerType == layers.LayerTypeLinuxSLL {
			// fmt.Println("Linux SLL Addr: ", sllLayer.Addr)
			// fmt.Println("Linux SLL Content: ", hex.EncodeToString(sllLayer.Contents))
			wireshark.Length = len(sllLayer.Payload)
			wireshark.SrcMac = sllLayer.Addr.String()
		}
		if layerType == layers.LayerTypeEthernet {
			// fmt.Println("Source MAC: ", ethLayer.SrcMAC)
			// fmt.Println("Destination MAC: ", ethLayer.DstMAC)
			// Ethernet type is typically IPv4 but could be ARP or other
			// fmt.Println("Ethernet type: ", ethLayer.EthernetType)
			wireshark.SrcMac = ethLayer.SrcMAC.String()
			wireshark.DstMac = ethLayer.DstMAC.String()
			wireshark.Length = len(ethLayer.Payload)
		}
		if layerType == layers.LayerTypeIPv4 {
			// fmt.Println("IPv4: ", ipLayer.SrcIP, "->", ipLayer.DstIP)
			// fmt.Println("IPv4 Content: ", hex.EncodeToString(ipLayer.Contents))
			wireshark.Source = ipLayer.SrcIP.String()
			wireshark.Destination = ipLayer.DstIP.String()
		}
		if layerType == layers.LayerTypeUDP {
			// fmt.Println("UDP Port: ", udpLayer.SrcPort, "->", udpLayer.DstPort)
			// fmt.Println("UDP Content: ", hex.EncodeToString(udpLayer.Contents))
			// fmt.Println("UDP Payload:", string(udpLayer.Payload))
			wireshark.SrcPort = udpLayer.SrcPort.String()
			wireshark.DstPort = udpLayer.DstPort.String()
		}
	}
	if err != nil {
		// fmt.Println("Trouble decoding layers: ", err)

		applicationLayer := (*packet).ApplicationLayer()
		if applicationLayer != nil {
			playload := applicationLayer.Payload()
			// fmt.Println("Playload Content: ", hex.EncodeToString(playload))
			// fmt.Println(playload)

			rtp := rtp.RTP{}
			err := rtp.Decode(playload)
			if err != nil {
				wireshark.Protocol = "UDP"
				wireshark.Info = fmt.Sprintf("%s -> %s Len(%v)", wireshark.SrcPort, wireshark.DstPort, len(playload))
			}
			wireshark.Protocol = "RTP"
			// wireshark.RtpPacket.SequenceNumber = rtp.SequenceNumber
			// wireshark.Timestamp = rtp.Timestamp
			wireshark.RtpPacket = &rtp
			wireshark.Info = rtp.ToWireSharkString()
		}
	}
	return &wireshark
}

func printAllPacketsInfo(rwp *rtp.WirePackets) {
	for _, pkt := range *rwp {
		fmt.Println(pkt.ToString())
	}
}

func printAllGroupPacketsInfo(rwpg *[]*rtp.WirePackets) {
	for _, rwp := range *rwpg {
		printAllPacketsInfo(rwp)
	}
}

func DoDecode(pcapFile, rtpSsrc, rtpCsrc string) (*rtp.WirePackets, *rtp.WirePackets) {
	// start := time.Now()

	// Open file instead of device
	handle, err = pcap.OpenOffline(pcapFile)
	if err != nil {
		// log.Fatal(err)
		log.Errorf("pcap.OpenOffline error, %v", err)
		return nil, nil
	}
	defer handle.Close()

	// Loop through packets in file
	// packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	var srcPacket []*rtp.WireShark
	var dstPacket []*rtp.WireShark

	for packet := range packetSource.Packets() {
		// wireshark := ParsePackets(packet)
		// faster
		// fmt.Println(packet)

		wireshark := DecodeLayers(&packet)
		if wireshark.RtpPacket == nil {
			continue
		}
		if rtpSsrc == "" && rtpCsrc == "" {
			if wireshark.SrcInPorts(udpPorts) {
				srcPacket = append(srcPacket, wireshark)
			} else if wireshark.DstInPorts(udpPorts) {
				dstPacket = append(dstPacket, wireshark)
			}
		}

		if rtpSsrc != "" && rtpCsrc == "" && wireshark.HasRtpSsrc(rtpSsrc) {
			if wireshark.SrcInPorts(udpPorts) {
				srcPacket = append(srcPacket, wireshark)
			} else if wireshark.DstInPorts(udpPorts) {
				dstPacket = append(dstPacket, wireshark)
			}
		}

		if rtpSsrc == "" && rtpCsrc != "" && wireshark.HasRtpCsrc(rtpCsrc) {
			if wireshark.SrcInPorts(udpPorts) {
				srcPacket = append(srcPacket, wireshark)
			} else if wireshark.DstInPorts(udpPorts) {
				dstPacket = append(dstPacket, wireshark)
			}
		}

		if rtpSsrc != "" && rtpCsrc != "" && wireshark.HasRtpSsrc(rtpSsrc) && wireshark.HasRtpCsrc(rtpCsrc) {
			if wireshark.SrcInPorts(udpPorts) {
				srcPacket = append(srcPacket, wireshark)
			} else if wireshark.DstInPorts(udpPorts) {
				dstPacket = append(dstPacket, wireshark)
			}
		}

		// printPacketInfo(packet)
	}
	// cost := time.Since(start)
	// fmt.Printf("执行时间: %v\n", cost)
	return (*rtp.WirePackets)(&dstPacket), (*rtp.WirePackets)(&srcPacket)
}
