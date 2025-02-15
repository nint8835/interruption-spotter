package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/nint8835/interruption-spotter/pkg/config"
	"github.com/nint8835/interruption-spotter/pkg/database"
	"github.com/nint8835/interruption-spotter/pkg/monitor"
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

		mon := monitor.New(db, cfg)

		srv := server.New(cfg, db)

		mon.Start()
		srv.Start()

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc

		srv.Stop(context.Background())
		mon.Stop()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
