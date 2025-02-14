package cmd

import (
	"github.com/spf13/cobra"

	"github.com/nint8835/interruption-spotter/pkg/config"
	"github.com/nint8835/interruption-spotter/pkg/database"
	"github.com/nint8835/interruption-spotter/pkg/server"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the app.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		checkErr(err, "Failed to load config")

		db, err := database.Connect(cfg)
		checkErr(err, "Failed to connect to database")

		srv, err := server.New(cfg, db)
		checkErr(err, "Failed to create server")

		err = srv.Run()
		checkErr(err, "Failed to run server")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
