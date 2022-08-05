package maintenance

import (
	"bytes"
	"html/template"
	"log"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/RocketChat/statuscentral/cmd/statusctl/common"
	"github.com/RocketChat/statuscentral/models"
)

var maintenanceDetailTemplate = `
ID: {{.ID}}
Title: {{.Title}}
Description: {{.Description}}
Created: {{ .CreatedAt.Format "Jan 02 2006 15:04" }}
Completed: {{ .Completed }}
Planned Start: {{ .PlannedStart.Format "Jan 02 2006 15:04" }}
Planned End: {{ .PlannedEnd.Format "Jan 02 2006 15:04" }}
Services: 
{{ range $service := .Services }}
- Name: {{$service.Name}}
  Status: {{$service.Status}}
  Regions:
  {{ range $region := $service.Regions }}
  - 
  {{ end }}
{{ end }}

Updates:
{{ range $update := .Updates }}
- ID: {{$update.ID}}
  Time: {{ $update.Time.Format "Jan 02 15:04" }}
  Status: {{ $update.Status }}
  Message: {{ $update.Message }}
{{ end }}
`

var describeCmd = &cobra.Command{
	Use:     "describe",
	Short:   "describe maintenance",
	Example: "statusctl maintenance describe",
	Run: func(c *cobra.Command, args []string) {
		client := common.GetStatusCentralClient()

		id, err := strconv.Atoi(args[0])
		if err != nil {
			panic("Unable to parse maintenance id")
		}

		maintenance, err := client.ScheduledMaintenance().Get(id)
		if err != nil {
			panic(err)
		}

		rendered, err := renderMaintenance(maintenance)
		if err != nil {
			panic(err)
		}

		log.Println(rendered)

	},
}

func renderMaintenance(scheduledMaintenance *models.ScheduledMaintenance) (string, error) {
	temp := template.New("maintenanceDetailTemplate")
	t, err := temp.Parse(maintenanceDetailTemplate)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, scheduledMaintenance); err != nil {
		return "", err
	}

	return buf.String(), nil
}
