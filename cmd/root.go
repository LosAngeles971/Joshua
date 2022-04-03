// CLI implementation
package cmd

import (
	"io/ioutil"
	"it/losangeles971/joshua/business/knowledge"
	"strings"
	
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var knowledgeFile string
var stateFile string
var success string
var maxCycles int

// CLI command to apply knowledge
var rootCmd = &cobra.Command{
	Use:   "joshua",
	Short: "joshua",
	Run: func(cmd *cobra.Command, args []string) {
		source, err := ioutil.ReadFile(knowledgeFile)
		if err != nil {
			log.Fatal(err)
		}
		data, err := ioutil.ReadFile(stateFile)
		if err != nil {
			log.Fatal(err)
		}
		var enc int
		if strings.HasSuffix(stateFile, ".json") {
			enc = knowledge.ENC_JSON
		} else {
			enc = knowledge.ENC_YAML
		}
		state := knowledge.NewState(knowledge.WithData(data, enc))
		engine, err := knowledge.NewEngine(string(source), maxCycles)
		if err != nil {
			log.Fatal(err)
		}
		solution := engine.IsItGoingToHappen(*state, success)
		if solution.Err != nil {
			log.Fatal(err)
		}
		solution.PrintChain()
		solution.PrintSummary()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&knowledgeFile, "knowledge", "k", "", "knowledge file")
	rootCmd.Flags().StringVarP(&stateFile, "data", "d", "", "initial state")
	rootCmd.Flags().StringVarP(&success, "success", "s", "", "final event")
	rootCmd.Flags().IntVarP(&maxCycles, "max-cycles", "m", 100, "maximum number of cycles (default 100)")
	rootCmd.MarkFlagRequired("knowledge")
	rootCmd.MarkFlagRequired("data")
	rootCmd.MarkFlagRequired("success")
}