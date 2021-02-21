package cmd

import (
	"fmt"
	"it/losangeles971/joshua/internal/knowledge"
	"os"
	"github.com/spf13/cobra"
)

var effectName string
// Command to execute a reason
var knowledgeCmd = &cobra.Command{
	Use:   "knowledge",
	Short: "Browse knowledge",
	Long: `Browse the knowledge.
Usage:
	joshua knowledge --knowledge|-k <knowledge-file> --effect <effect>`,
	Run: func(cmd *cobra.Command, args []string) {
		k := knowledge.Knowledge{}
		err := k.Load(knowledgeFile)
		if err != nil {
			fmt.Println("Error loading the problem: ", err)
			os.Exit(1)
		}
		if len(effectName) > 0 {
			_, ok := k.GetEvent(effectName)
			if !ok {
				fmt.Println("Effect not existent: ", effectName)
				os.Exit(1)
			}
		}
	},
}

func init() {
	knowledgeCmd.Flags().StringVarP(&knowledgeFile, "knowledge", "k", "", "YAML file representing the knowledge")
	knowledgeCmd.Flags().StringVarP(&effectName, "effect", "e", "", "Success effect")
	solveCmd.MarkFlagRequired("knowledge")
	rootCmd.AddCommand(knowledgeCmd)
}


