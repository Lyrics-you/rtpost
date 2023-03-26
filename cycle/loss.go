package cycle

import (
	"fmt"
	"rtpost/rtp"
)

func printSingleLossInfo(rwp *rtp.WirePackets) {
	if len(*rwp) == 0 {
		return
	}
	loss, info := rwp.LostSequenceNumbers()
	fmt.Printf("SSRC=%s total loss : %d\n", (*rwp)[0].RtpPacket.ToSSRCString(), loss)
	for _, i := range *info {
		fmt.Println(i)
	}
}

func showLossInfo(rwp *rtp.WirePackets, group bool) {
	if len(*rwp) == 0 {
		return
	}
	rwp = rwp.RemoveNotTargetType()
	if group {
		dGroup := rwp.DoGroup()
		for _, dg := range *dGroup {
			printSingleLossInfo(dg)
		}
	} else {
		printSingleLossInfo(rwp)
	}
}

func DoLoss(pcapFile, rtpSsrc, rtpCsrc string, group bool) {
	if rtpSsrc == "" && rtpCsrc == "" {
		group = true
	}
	dstPacket, srcPacket := DoDecode(pcapFile, rtpSsrc, rtpCsrc)
	showLossInfo(dstPacket, group)
	showLossInfo(srcPacket, group)
}
