package cmd

import (
	"errors"
	"fmt"
	"rtpost/logger"
	"rtpost/version"
	"strings"

	"github.com/spf13/cobra"
)

var (
	log           = logger.Logger()
	filePath      string
	rtpSsrc       string
	rtpCsrc       string
	deduplicate   bool
	group         bool
	delta         uint32
	rtpostVersion bool
)

var rootCmd = &cobra.Command{
	Use:   "rtpost",
	Short: "rtpost is a analyze rtp packet loss rate system.",
	Long: `
rtpost is a free and open source analyze live rtp packet loss rate system,
designed to analyse the udp packets of the specified port in the pcap file obtained by tcpdump, 
decode them into rtp packets and then analyse them to obtain the packet loss rate.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			Error(cmd, args, errors.New("unrecognized parameters"))
			return
		}
		if rtpostVersion {
			showVersion()
			return
		}

	},
}

func upperSSRC(rtpSsrc string) string {
	if strings.HasPrefix(rtpSsrc, "0x") {
		return "0x" + strings.ToUpper(rtpSsrc[2:])
	}
	return rtpSsrc
}

func showVersion() {
	historys := version.Historys
	history := historys[len(historys)-1]
	fmt.Printf("rtpost(%s):%s\n", history.Version, history.Description)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "pcap file path")
	rootCmd.PersistentFlags().StringVarP(&rtpSsrc, "ssrc", "s", "", "rtp ssrc is hexadecimal, prefixed with 0x")
	rootCmd.PersistentFlags().StringVarP(&rtpCsrc, "csrc", "c", "", "rtp csrc is hexadecimal, prefixed with 0x")
	rootCmd.PersistentFlags().BoolVarP(&deduplicate, "deduplicate", "d", false, "rtp info deduplicate by ip, port and ssrc")
	rootCmd.PersistentFlags().BoolVarP(&group, "group", "g", false, "rtp info deduplicate in group")
	rootCmd.PersistentFlags().BoolVarP(&rtpostVersion, "version", "v", false, "show rtpost version")
	rootCmd.PersistentFlags().Uint32VarP(&delta, "time-delta", "t", 0, "show bigger than delta's packet")
}
func Execute() {
	rootCmd.Execute()
}
