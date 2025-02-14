package cmd

import (
	"context"
	"database/sql"
	"errors"

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
					currentLevel, err := db.GetCurrentInterruptionLevel(ctx, database.GetCurrentInterruptionLevelParams{
						Region:          regionName,
						OperatingSystem: osName,
						InstanceType:    instanceType,
					})
					if !errors.Is(err, sql.ErrNoRows) {
						checkErr(err, "Failed to get current interruption level")
					}

					if currentLevel == int64(instanceStats.InterruptionLevel) {
						continue
					}

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
