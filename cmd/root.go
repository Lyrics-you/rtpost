package cmd

import (
	"errors"
	"rtpost/logger"
	"rtpost/stat"

	"github.com/spf13/cobra"
)

var (
	log      = logger.Logger()
	filePath string
	rtpSsrc  string
	visual   bool
)

var rootCmd = &cobra.Command{
	Use:   "rtpost",
	Short: "rtpost is a analyze media rtp packet loss rate system.",
	Long: `
rtpost is a free and open source analyze media rtp packet loss rate system,
designed to analyse the udp packets of the specified port in the pcap file obtained by tcpdump, 
decode them into rtp packets and then analyse them to obtain the packet loss rate.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			Error(cmd, args, errors.New("unrecognized parameters"))
			return
		}
		if filePath != "" && rtpSsrc != "" {
			stat.DoStat(filePath, rtpSsrc, visual)
		} else {
			Error(cmd, args, errors.New("invalid parameters"))
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "pcap file path")
	rootCmd.PersistentFlags().StringVarP(&rtpSsrc, "ssrc", "s", "", "rtp ssrc")
	rootCmd.PersistentFlags().BoolVarP(&visual, "visual", "v", false, "rtp info visualization")
}
func Execute() {
	rootCmd.Execute()
}
