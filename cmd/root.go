package cmd

import (
	"io/ioutil"
	"it/losangeles971/joshua/engine"
	"it/losangeles971/joshua/knowledge"
	"it/losangeles971/joshua/outputs"
	"log"

	"github.com/spf13/cobra"
)

var knowledgeFile string
var problemFile string
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
		k, err := knowledge.Load(string(source))
		if err != nil {
			log.Fatal(err)
		}
		source, err = ioutil.ReadFile(problemFile)
		if err != nil {
			log.Fatal(err)
		}
		s, success_name, err := engine.LoadProblem(string(source))
		if err != nil {
			log.Fatal(err)
		}
		success, ok := k.GetEvent(success_name)
		if !ok {
			log.Fatalf("there is not event '%v' into the knowledge", success_name)
		}
		outcome, queue, err := engine.MakeItHappen(k, *s, success, maxCycles)
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
	rootCmd.Flags().StringVarP(&problemFile, "problem", "p", "", "YAML file describing the problem")
	rootCmd.Flags().StringVarP(&solutionFile, "output", "o", "", "YAML file to host the solution")
	rootCmd.Flags().StringVarP(&dotFile, "dot", "d", "", "Graphwiz Dot file of the solution")
	rootCmd.Flags().IntVarP(&maxCycles, "max-cycles", "m", 100, "Maximum number of cycles (default 100)")
	rootCmd.MarkFlagRequired("knowledge")
	rootCmd.MarkFlagRequired("problem")
}