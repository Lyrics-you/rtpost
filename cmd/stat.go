package cmd

import (
	"errors"
	"rtpost/cycle"

	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "stat rtp packets's information",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			Error(cmd, args, errors.New("unrecognized command"))
			return
		}
		if filePath != "" {
			cycle.DoStat(filePath, upperSSRC(rtpSsrc), upperSSRC(rtpCsrc), deduplicate, group)
		} else {
			Error(cmd, args, errors.New("invalid parameters"))
		}
	},
}

func init() {
	rootCmd.AddCommand(statCmd)
}
