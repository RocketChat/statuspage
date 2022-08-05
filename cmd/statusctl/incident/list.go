package incident

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	"github.com/RocketChat/statuscentral/cmd/statusctl/common"
)

var latestOnly = false

var listCmd = &cobra.Command{
	Use: "list",
	Aliases: []string{
		"ls",
	},
	Short:   "List incidents",
	Example: "statusctl incidents ls",
	Run: func(c *cobra.Command, args []string) {
		t := table.NewWriter()

		t.Style().Options.DrawBorder = false
		t.Style().Options.SeparateRows = false
		t.Style().Options.SeparateColumns = false
		t.Style().Options.SeparateHeader = false
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ID", "Title", "Status", "Date"})

		cl := common.GetStatusCentralClient()

		incidents, err := cl.Incidents().GetMultiple(latestOnly)
		if err != nil {
			panic(err)
		}

		for _, incident := range incidents {
			t.AppendRows([]table.Row{
				{incident.ID, incident.Title, incident.Status.String(), incident.Time.Format("Jan 02 2006 15:04")},
			})
		}

		t.Render()

	},
}
