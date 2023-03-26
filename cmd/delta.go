package cmd

import (
	"errors"
	"rtpost/cycle"

	"github.com/spf13/cobra"
)

var deltaCmd = &cobra.Command{
	Use:   "delta",
	Short: "delta rtp packets",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			Error(cmd, args, errors.New("unrecognized command"))
			return
		}
		if filePath != "" {
			cycle.DoDelta(filePath, upperSSRC(rtpSsrc), upperSSRC(rtpCsrc), delta, deduplicate, group)
		} else {
			Error(cmd, args, errors.New("invalid parameters"))
		}
	},
}

func init() {
	rootCmd.AddCommand(deltaCmd)
}
