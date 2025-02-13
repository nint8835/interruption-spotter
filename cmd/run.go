package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nint8835/interruption-spotter/pkg/config"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the app.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		checkErr(err, "Failed to load config")

		fmt.Printf("%#+v\n", cfg)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
