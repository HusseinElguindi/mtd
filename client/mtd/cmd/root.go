package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mtd",
	Short: "mtd stands for MultiThreaded Downloader",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("please use a command")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
