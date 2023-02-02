package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

func Error(cmd *cobra.Command, args []string, err error) {
	log.Errorf("execute '%s %s' error, %v", cmd.Name(), strings.Join(args, " "), err)
}
