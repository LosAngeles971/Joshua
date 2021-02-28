package cmd

import (
	"fmt"
	"it/losangeles971/joshua/internal/knowledge"
	"it/losangeles971/joshua/internal/problems"
	"it/losangeles971/joshua/internal/outputs"
	"it/losangeles971/joshua/pkg"
	"os"
	"github.com/spf13/cobra"
)

var knowledgeFile string
var problemFile string
var solutionFile string
var maxCycles int
var dotFile string

// Command to execute a reason
var solveCmd = &cobra.Command{
	Use:   "solve",
	Short: "Solve a problem",
	Long: `Solve a problem using a cause-effect approach.
Usage:
	joshua solve --knowledge|-k <knowledge-file> --problem|-p <problem-file>`,
	Run: func(cmd *cobra.Command, args []string) {
		k := knowledge.Knowledge{}
		err := k.Load(knowledgeFile)
		if err != nil {
			fmt.Println("Error loading the knowledge: ", err)
			os.Exit(1)
		}
		init, s_name, err := problems.Load(problemFile)
		if err != nil {
			fmt.Println("Error loading the problem: ", err)
			os.Exit(1)
		}
		success, ok := k.GetEvent(s_name)
		if !ok {
			fmt.Println("Problem declares an event not included by knowledge: ", s_name)
			os.Exit(1)
		}
		outcome, queue, err := pkg.Reason(k, init, success, maxCycles)
		if err != nil {
			fmt.Println("Error solving the problem: ", err)
			os.Exit(1)
		}
		pkg.PrintSummary(outcome, queue)
		if len(solutionFile) > 0 {
			err = queue.Save(solutionFile)
			if err != nil {
				fmt.Println("Error writing the solution file: ", err)
				os.Exit(1)
			}
		}
		if len(dotFile) > 0 {
			err = outputs.SaveDot(queue, dotFile)
			if err != nil {
				fmt.Println("Error writing the dot solution file: ", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	solveCmd.Flags().StringVarP(&knowledgeFile, "knowledge", "k", "", "YAML file representing the knowledge")
	solveCmd.Flags().StringVarP(&problemFile, "problem", "p", "", "YAML file describing the problem")
	solveCmd.Flags().StringVarP(&solutionFile, "output", "o", "", "YAML file to host the solution")
	solveCmd.Flags().StringVarP(&dotFile, "dot", "d", "", "Graphwiz Dot file of the solution")
	solveCmd.Flags().IntVarP(&maxCycles, "max-cycles", "m", 100, "Maximum number of cycles (default 100)")
	solveCmd.MarkFlagRequired("knowledge")
	solveCmd.MarkFlagRequired("problem")
	rootCmd.AddCommand(solveCmd)
}


