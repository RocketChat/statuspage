package incident

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
	Short:   "get incident",
	Example: "statusctl incident get",
	Run: func(c *cobra.Command, args []string) {
		client := common.GetStatusCentralClient()

		id, err := strconv.Atoi(args[0])
		if err != nil {
			panic("Unable to parse incident id")
			return
		}

		incident, err := client.Incidents().Get(id)
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
			t.AppendHeader(table.Row{"ID", "Title", "Status", "Date"})
			t.AppendRows([]table.Row{
				{incident.ID, incident.Title, incident.Status.String(), incident.Time.Format("Jan 02 2006 15:04")},
			})
			t.Render()
		case "json":
			jsonText, err := json.Marshal(incident)
			if err != nil {
				panic(err)
			}

			log.Println(string(jsonText))
		}

	},
}
