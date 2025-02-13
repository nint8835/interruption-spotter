package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nint8835/interruption-spotter/pkg/config"
	"github.com/nint8835/interruption-spotter/pkg/spotdata"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the app.",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := config.Load()
		checkErr(err, "Failed to load config")

		fetcher := spotdata.Fetcher{}

		shouldFetch, err := fetcher.ShouldFetch(context.Background())
		checkErr(err, "Failed to check if we should fetch")

		fmt.Println(shouldFetch)

		resp, err := fetcher.Fetch(context.Background())
		checkErr(err, "Failed to fetch data")

		fmt.Printf("%#+v\n", resp)

		nowShouldFetch, err := fetcher.ShouldFetch(context.Background())
		checkErr(err, "Failed to check if we should fetch")

		fmt.Println(nowShouldFetch)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
