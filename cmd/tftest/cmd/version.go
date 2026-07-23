package cmd

import (
	"fmt"

	"github.com/matthew-dresden/terraform-terratest-framework/cmd/tftest/logger"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("TFTest CLI %s", Version)
		fmt.Printf("🎉 TFTest CLI %s 🎉\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
