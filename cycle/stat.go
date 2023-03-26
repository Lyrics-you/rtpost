package cycle

import (
	"fmt"
	"rtpost/rtp"
)

func LossRate(loss, expected int32) float64 {
	return float64(loss) / float64(expected)
}

func StatInfo(wp *rtp.WirePackets, recv, expect, loss int32) string {
	ws := (*wp)[0]
	return fmt.Sprintf("%s:%s->%s:%s {SSRC:%s CSRC:%s} (recv:%d, expect:%d, loss:%d, rate:%.2f%%)", ws.Source, ws.SrcPort, ws.Destination, ws.DstPort, ws.RtpPacket.ToSSRCString(), ws.RtpPacket.ToCSRCItemString(), recv, expect, loss, LossRate(loss, expect)*100)
}

func printSingleStatInfo(rwp *rtp.WirePackets, deduplicate bool) {
	if len(*rwp) == 0 {
		return
	}
	expected, loss := rwp.ExpectedPacketsLengthAndLoss()
	if deduplicate {
		fmt.Println(StatInfo(rwp, expected-loss, expected, loss))
	} else {
		fmt.Println(StatInfo(rwp, int32(len(*rwp)), expected, expected-int32(len(*rwp))))
	}
}

func showStatInfo(rwp *rtp.WirePackets, deduplicate, group bool) {
	if len(*rwp) == 0 {
		return
	}
	rwp = rwp.RemoveNotTargetType()
	if group {
		dGroup := rwp.DoGroup()
		for _, dg := range *dGroup {
			printSingleStatInfo(dg, deduplicate)
		}
	} else {
		printSingleStatInfo(rwp, deduplicate)
	}
}

func DoStat(pcapFile, rtpSsrc, rtpCsrc string, deduplicate, group bool) {
	if rtpSsrc == "" && rtpCsrc == "" {
		group = true
	}
	dstPacket, srcPacket := DoDecode(pcapFile, rtpSsrc, rtpCsrc)
	showStatInfo(dstPacket, deduplicate, group)
	showStatInfo(srcPacket, deduplicate, group)
}
