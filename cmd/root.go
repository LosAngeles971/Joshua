// CLI implementation
package cmd

import (
	"io/ioutil"
	"it/losangeles971/joshua/engine"
	"it/losangeles971/joshua/knowledge"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var knowledgeFile string
var dataFile string
var success_name string
var solutionFile string
var maxCycles int

// CLI command to apply knowledge
var rootCmd = &cobra.Command{
	Use:   "joshua",
	Short: "joshua",
	Long:  `joshua::
	Applying knowledge to understand if a "success" may occur starting from a given initial state.
`,
	Run: func(cmd *cobra.Command, args []string) {
		source, err := ioutil.ReadFile(knowledgeFile)
		if err != nil {
			log.Fatal(err)
		}
		s := knowledge.NewState(knowledge.WithYAML(dataFile))
		outcome, queue, err := engine.IsItGoingToHappen(string(source), *s, success_name, maxCycles)
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
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&knowledgeFile, "knowledge", "k", "", "YAML file representing the knowledge")
	rootCmd.Flags().StringVarP(&dataFile, "data", "d", "", "YAML file describing the problem")
	rootCmd.Flags().StringVarP(&solutionFile, "output", "o", "", "YAML file to host the solution")
	rootCmd.Flags().StringVarP(&success_name, "success", "s", "", "Name of success event")
	rootCmd.Flags().IntVarP(&maxCycles, "max-cycles", "m", 100, "Maximum number of cycles (default 100)")
	rootCmd.MarkFlagRequired("knowledge")
	rootCmd.MarkFlagRequired("problem")
}