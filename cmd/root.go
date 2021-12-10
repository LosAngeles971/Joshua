/*
definition of the command line flags
*/
package cmd

import (
	"io/ioutil"
	"it/losangeles971/joshua/engine"
	"it/losangeles971/joshua/outputs"
	"log"
	"os"
	
	"github.com/spf13/cobra"
)

var knowledgeFile string
var dataFile string
var success_name string
var solutionFile string
var maxCycles int
var dotFile string

var rootCmd = &cobra.Command{
	Use:   "joshua",
	Short: "joshua",
	Long:  `joshua
	Usage:
		joshua --knowledge|-k <knowledge-file> --problem|-p <problem-file>`,
	Run: func(cmd *cobra.Command, args []string) {
		source, err := ioutil.ReadFile(knowledgeFile)
		if err != nil {
			log.Fatal(err)
		}
		outcome, queue, err := engine.IsItGoingToHappen(string(source), nil, success_name, maxCycles)
		if err != nil {
			log.Fatal(err)
		}
		engine.PrintSummary(outcome, queue)
		if len(solutionFile) > 0 {
			err = queue.Save(solutionFile)
			if err != nil {
				log.Fatal(err)
			}
		}
		if len(dotFile) > 0 {
			err = outputs.SaveDot(queue, dotFile)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&knowledgeFile, "knowledge", "k", "", "YAML file representing the knowledge")
	rootCmd.Flags().StringVarP(&dataFile, "data", "d", "", "YAML file describing the problem")
	rootCmd.Flags().StringVarP(&solutionFile, "output", "o", "", "YAML file to host the solution")
	rootCmd.Flags().StringVarP(&dotFile, "dot", "d", "", "Graphwiz Dot file of the solution")
	rootCmd.Flags().StringVarP(&success_name, "success", "s", "", "Name of success event")
	rootCmd.Flags().IntVarP(&maxCycles, "max-cycles", "m", 100, "Maximum number of cycles (default 100)")
	rootCmd.MarkFlagRequired("knowledge")
	rootCmd.MarkFlagRequired("problem")
}