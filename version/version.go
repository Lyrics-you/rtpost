package version

type History struct {
	Version     string
	Description string
}

var (
	Historys = []History{
		{Version: "0.1.0",
			Description: "first version, a analyze rtp packet loss rate system",
		},

		{Version: "0.2.0",
			Description: "add deduplication feature",
		}, {
			Version:     "0.2.1",
			Description: "change rtp length to uint32",
		},

		{Version: "0.3.0",
			Description: "adjust RTP expected packet length and RTP loss rate",
		}, {
			Version:     "0.3.1",
			Description: "rtp ssrc supports lower case",
		}, {
			Version:     "0.3.2",
			Description: "rtpost supports --version to see version and description",
		}, {
			Version:     "0.3.3",
			Description: "deduplicat rtp packets will remove not dynamicRTP",
		},

		{Version: "0.4.0",
			Description: "no -s will seek all SSRC, now -g will group by SSRC",
		}, {
			Version:     "0.4.1",
			Description: "fix deduplicat rtp packets will remove not dynamicRTP-116, change default time delat to 27000",
		}, {
			Version:     "0.4.2",
			Description: "rtpost supports -c to display a particular CSRC",
		}, {
			Version:     "0.4.3",
			Description: "update HasRtpSSRC and HasRtpCSRC function",
		}, {
			Version:     "0.4.4",
			Description: "fix wireshark.RtpPacket nil's problem",
		},

		{Version: "0.5.0",
			Description: "performance improvement",
		}, {
			Version:     "0.5.1",
			Description: "fix gourp's key to src:srcport->dst:dstport-ssrc",
		}, {
			Version:     "0.5.2",
			Description: "fix ExpectedPacketsLength function",
		}, {
			Version:     "0.5.3",
			Description: "fix filter out unnecessary packets in RTP",
		},

		{Version: "0.6.0",
			Description: "new version : offer decode, stat, loss, delta command ",
		},
	}
)
