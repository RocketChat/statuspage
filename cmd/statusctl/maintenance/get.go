package maintenance

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"

	"github.com/RocketChat/statuscentral/cmd/statusctl/common"
)

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "get maintenance",
	Example: "statusctl maintenance get",
	Run: func(c *cobra.Command, args []string) {
		client := common.GetStatusCentralClient()

		id, err := strconv.Atoi(args[0])
		if err != nil {
			panic("Unable to parse maintenance id")
			return
		}

		maintenance, err := client.ScheduledMaintenance().Get(id)
		if err != nil {
			panic(err)
		}

		switch outputFormat {
		case "list":
			t := table.NewWriter()

			t.Style().Options.DrawBorder = false
			t.Style().Options.SeparateRows = false
			t.Style().Options.SeparateColumns = false
			t.Style().Options.SeparateHeader = false
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"ID", "Title", "Planned Start", "Planned End", "Completed"})
			t.AppendRows([]table.Row{
				{maintenance.ID, maintenance.Title, maintenance.PlannedStart.Format("Jan 02 2006 15:04"), maintenance.PlannedEnd.Format("Jan 02 2006 15:04"), maintenance.Completed},
			})
			t.Render()
		case "json":
			jsonText, err := json.Marshal(maintenance)
			if err != nil {
				panic(err)
			}

			log.Println(string(jsonText))
		}

	},
}
