package cycle

import "rtpost/rtp"

func DoDeduplicate(pcapFile, rtpSsrc, rtpCsrc string, deduplicate, group bool) {
	if rtpSsrc == "" && rtpCsrc == "" {
		group = true
	}
	dstPacket, srcPacket := DoDecode(pcapFile, rtpSsrc, rtpCsrc)
	dstGroup := &[]*rtp.WirePackets{}
	srcGroup := &[]*rtp.WirePackets{}
	if deduplicate {
		if group {
			dstGroup = dstPacket.DeduplicateGroup()
			srcGroup = srcPacket.DeduplicateGroup()
		} else {
			dstPacket = dstPacket.Deduplicate()
			srcPacket = srcPacket.Deduplicate()
		}
	} else {
		if group {
			dstGroup = dstPacket.DoGroup()
			srcGroup = srcPacket.DoGroup()
		}
	}
	if group {
		printAllGroupPacketsInfo(dstGroup)
		printAllGroupPacketsInfo(srcGroup)
	} else {
		printAllPacketsInfo(dstPacket)
		printAllPacketsInfo(srcPacket)
	}
}
