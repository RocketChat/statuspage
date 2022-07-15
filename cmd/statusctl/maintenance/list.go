package maintenance

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
	Short:   "List scheduled maintenance",
	Example: "statusctl maintenance ls",
	Run: func(c *cobra.Command, args []string) {
		t := table.NewWriter()

		t.Style().Options.DrawBorder = false
		t.Style().Options.SeparateRows = false
		t.Style().Options.SeparateColumns = false
		t.Style().Options.SeparateHeader = false
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ID", "Title", "Planned Start", "Planned End", "Completed"})

		cl := common.GetStatusCentralClient()

		maintenances, err := cl.ScheduledMaintenance().GetMultiple(latestOnly)
		if err != nil {
			panic(err)
		}

		for _, maintenance := range maintenances {
			t.AppendRows([]table.Row{
				{maintenance.ID, maintenance.Title, maintenance.PlannedStart.Format("Jan 02 2006 15:04"), maintenance.PlannedEnd.Format("Jan 02 2006 15:04"), maintenance.Completed},
			})
		}

		t.Render()

	},
}
