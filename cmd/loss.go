package cmd

import (
	"errors"
	"rtpost/cycle"

	"github.com/spf13/cobra"
)

var lossCmd = &cobra.Command{
	Use:   "loss",
	Short: "loss show rtp lost packets",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			Error(cmd, args, errors.New("unrecognized command"))
			return
		}
		if filePath != "" {
			cycle.DoLoss(filePath, upperSSRC(rtpSsrc), upperSSRC(rtpCsrc), group)
		} else {
			Error(cmd, args, errors.New("invalid parameters"))
		}
	},
}

func init() {
	rootCmd.AddCommand(lossCmd)
}
