package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "interruption-spotter",
	Short: "RSS feed alerting of changes to AWS Spot Instance interruption rates.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		// TODO: Don't panic
		panic(err)
	}
}
