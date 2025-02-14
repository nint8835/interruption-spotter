package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/nint8835/interruption-spotter/pkg/config"
	"github.com/nint8835/interruption-spotter/pkg/database"
	"github.com/nint8835/interruption-spotter/pkg/spotdata"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the app.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		checkErr(err, "Failed to load config")

		db, err := database.Connect(cfg)
		checkErr(err, "Failed to connect to database")

		fetcher := spotdata.Fetcher{}

		ctx := context.Background()

		resp, err := fetcher.Fetch(ctx)
		checkErr(err, "Failed to fetch data")

		for regionName, regionStats := range resp.SpotAdvisor {
			for osName, osStats := range regionStats {
				for instanceType, instanceStats := range osStats {
					err = db.InsertStat(ctx, database.InsertStatParams{
						Region:            regionName,
						OperatingSystem:   osName,
						InstanceType:      instanceType,
						InterruptionLevel: int64(instanceStats.InterruptionLevel),
					})
					checkErr(err, "Failed to insert stat")
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
