package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "shedock",
	Short: "shedock is a CLI application",
	Long:  `shedock is a CLI application. It accepts a filepath as an argument.`,
	Args:  cobra.ExactArgs(1),
	Run:   DriverCli,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of shedock",
	Long:  `All software has versions. This is shedock's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("shedock v1.0.0")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
