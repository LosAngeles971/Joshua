package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "joshua",
	Short: "Joshua is a Reasoner",
	Long:  `Joshua is a reasoner based on the principle of cause-effect,
which leverages on the graph to represent knowledge and applying reasoning.
2021 - LosAngeles971`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
