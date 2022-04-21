package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func version(cmd *cobra.Command) {

	var version = &cobra.Command{
		Use:          "version",
		Long:         `version`,
		Short:        "auto version",
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("version v.1.1")
		}}
	cmd.AddCommand(version)
}
