package pkg

import (
	kkk "it/losangeles971/joshua/internal/knowledge"
	ctx "it/losangeles971/joshua/internal/context"
	"text/tabwriter"
	"fmt"
	"os"
)

func PrintSummary(outcome string, queue kkk.Queue) {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Outcome\t" + outcome + "\t")
	fmt.Fprintln(w, "Cycles\t" + fmt.Sprint(queue.GetCycles()) + "\t")
	fmt.Fprintln(w, "Queue's size\t" + fmt.Sprint(queue.Size()) + "\t")
	w.Flush()
}

func Reason(k kkk.Knowledge, init ctx.State, effect *kkk.Event, max_cycles int) (string, kkk.Queue, error) {
	return kkk.Reason(k, init, effect, max_cycles)
}