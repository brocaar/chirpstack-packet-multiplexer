package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the ChirpStack Packet Multiplexer version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
