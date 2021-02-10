package cmd

import (
	"fmt"
	"it/losangeles971/joshua/internal/io"
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
		k, err := io.Load(knowledgeFile)
		if err != nil {
			fmt.Println("Error loading the problem: ", err)
			os.Exit(1)
		}
		if len(effectName) > 0 {
			effect, ok := k.GetEvent(effectName)
			if !ok {
				fmt.Println("Effect not existent: ", effectName)
				os.Exit(1)
			}
			for _, e := range k.IsEffectOf(effect) {
				fmt.Println(e.Print())
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


