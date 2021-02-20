package cmd

import (
	"fmt"
	"it/losangeles971/joshua/internal/io"
	"it/losangeles971/joshua/pkg"
	"os"
	"github.com/spf13/cobra"
)

var knowledgeFile string
var problemFile string
var solutionFile string
var maxCycles int
var queuelog bool

// Command to execute a reason
var solveCmd = &cobra.Command{
	Use:   "solve",
	Short: "Solve a problem",
	Long: `Solve a problem using a cause-effect approach.
Usage:
	joshua solve --knowledge|-k <knowledge-file> --problem|-p <problem-file>`,
	Run: func(cmd *cobra.Command, args []string) {
		init, k, success, err := io.LoadProblem(knowledgeFile, problemFile)
		if err != nil {
			fmt.Println("Error loading the problem: ", err)
			os.Exit(1)
		}
		outcome, queue, err := pkg.Reason(k, init, &success, maxCycles)
		if err != nil {
			fmt.Println("Error solving the problem: ", err)
			os.Exit(1)
		}	
		pkg.PrintSummary(outcome, queue)
	},
}

func init() {
	solveCmd.Flags().StringVarP(&knowledgeFile, "knowledge", "k", "", "YAML file representing the knowledge")
	solveCmd.Flags().StringVarP(&problemFile, "problem", "p", "", "YAML file describing the problem")
	solveCmd.Flags().StringVarP(&solutionFile, "output", "o", "", "YAML file to host the solution")
	solveCmd.Flags().IntVarP(&maxCycles, "max-cycles", "m", 100, "Maximum number of cycles (default 100)")
	solveCmd.MarkFlagRequired("knowledge")
	solveCmd.MarkFlagRequired("problem")
	solveCmd.MarkFlagRequired("output")
	rootCmd.AddCommand(solveCmd)
}


