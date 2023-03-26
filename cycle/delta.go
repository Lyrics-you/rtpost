package cycle

import (
	"fmt"
	"rtpost/rtp"
)

func showSingleDeltaInfo(rwp *rtp.WirePackets, delta uint32) {
	if len(*rwp) == 0 {
		return
	}
	rwp = rwp.Delta()
	deltaed := []string{}
	for _, i := range *rwp {
		if i.Delta.Seconds() >= float64(delta) {
			deltaed = append(deltaed, i.ToHasDeltaString())
		}
	}
	fmt.Printf("SSRC=%s delta(>=%ds)  : % d\n", (*rwp)[0].RtpPacket.ToSSRCString(), delta, len(deltaed))
	for _, i := range deltaed {
		fmt.Println(i)
	}

}

func showDeltaInfo(rwp *rtp.WirePackets, delta uint32, deduplicate, group bool) {
	if len(*rwp) == 0 {
		return
	}
	rwp = rwp.RemoveNotTargetType()
	if deduplicate {
		if group {
			dGroup := rwp.DeduplicateGroup()
			for _, dg := range *dGroup {
				showSingleDeltaInfo(dg, delta)
			}
		} else {
			showSingleDeltaInfo(rwp, delta)
		}
	} else {
		if group {
			dGroup := rwp.DoGroup()
			for _, dg := range *dGroup {
				showSingleDeltaInfo(dg, delta)
			}
		} else {
			showSingleDeltaInfo(rwp, delta)
		}
	}
}

func DoDelta(pcapFile, rtpSsrc, rtpCsrc string, delta uint32, deduplicate, group bool) {
	if rtpSsrc == "" && rtpCsrc == "" {
		group = true
	}
	dstPacket, srcPacket := DoDecode(pcapFile, rtpSsrc, rtpCsrc)

	showDeltaInfo(dstPacket, delta, deduplicate, group)
	showDeltaInfo(srcPacket, delta, deduplicate, group)
}

// func DoDeduplicate(pcapFile, rtpSsrc, rtpCsrc string, deduplicate, group bool) {
// 	if rtpSsrc == "" && rtpCsrc == "" {
// 		group = true
// 	}
// 	dstPacket, srcPacket := DoDecode(pcapFile, rtpSsrc, rtpCsrc)
// 	dstGroup := &[]*rtp.WirePackets{}
// 	srcGroup := &[]*rtp.WirePackets{}
// 	if deduplicate {
// 		if group {
// 			dstGroup = dstPacket.DeduplicateGroup()
// 			srcGroup = srcPacket.DeduplicateGroup()
// 		} else {
// 			dstPacket = dstPacket.Deduplicate()
// 			srcPacket = srcPacket.Deduplicate()
// 		}
// 	} else {
// 		if group {
// 			dstGroup = dstPacket.DoGroup()
// 			srcGroup = srcPacket.DoGroup()
// 		}
// 	}
// 	if group {
// 		printAllGroupPacketsInfo(dstGroup)
// 		printAllGroupPacketsInfo(srcGroup)
// 	} else {
// 		printAllPacketsInfo(dstPacket)
// 		printAllPacketsInfo(srcPacket)
// 	}
// }
