package maintenance

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/RocketChat/statuscentral/cmd/statusctl/common"
	"github.com/RocketChat/statuscentral/models"
	"github.com/spf13/cobra"
)

var patchCmd = &cobra.Command{
	Use:     "patch",
	Short:   "patch an existing maintenance",
	Example: "statusctl maintenance patch [id]",
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

		rendered, err := renderMaintenance(maintenance)
		if err != nil {
			panic(err)
		}

		log.Println(rendered)

		title := common.StringPrompt(fmt.Sprintf("Title [%s]:", maintenance.Title))
		if title == "" {
			title = maintenance.Title
		}

		description := common.StringPrompt(fmt.Sprintf("Description [%s]:", maintenance.Description))
		if description == "" {
			description = maintenance.Description
		}

		var plannedStartTime time.Time
		plannedStartTimeText := common.StringPrompt(fmt.Sprintf("Planned Start UTC Time [%s]:", maintenance.PlannedStart.Format("2006/01/02 15:04:05")))
		if plannedStartTimeText == "" {
			plannedStartTime = maintenance.PlannedStart
		} else {
			plannedStartTime, err = time.ParseInLocation("2006/01/02 15:04:05", plannedStartTimeText, time.UTC)
			if err != nil {
				panic(err)
			}
		}

		var plannedEndTime time.Time
		plannedEndTimeText := common.StringPrompt(fmt.Sprintf("Planned End UTC Time [%s]:", maintenance.PlannedEnd.Format("2006/01/02 15:04:05")))
		if plannedEndTimeText == "" {
			plannedEndTime = maintenance.PlannedEnd
		} else {
			plannedEndTime, err = time.ParseInLocation("2006/01/02 15:04:05", plannedEndTimeText, time.UTC)
			if err != nil {
				panic(err)
			}
		}

		updateMessage := common.StringPrompt("Status Update Message:")

		for i, statusOption := range models.IncidentStatusArray {
			log.Printf("%d) %s\n", i, statusOption)
		}

		status, err := common.IntPrompt("Updated Maintenance Status [3]:", 3)
		if err != nil {
			log.Fatalln("Invalid selection")
		}

		updateServiceStatus, err := common.GetYesNoPrompt("Update Service Status?", false)
		if err != nil {
			panic(err)
		}

		statusUpdate := &models.StatusUpdate{
			Status:   models.IncidentStatusArray[status],
			Message:  updateMessage,
			Services: maintenance.Services,
		}

		if updateServiceStatus {
			serviceUpdates, err := updateImpactedServices(maintenance.Services)
			if err != nil {
				panic(err)
			}

			statusUpdate.Services = serviceUpdates
		}

		maintenance.Title = title
		maintenance.Description = description
		maintenance.PlannedStart = plannedStartTime
		maintenance.PlannedEnd = plannedEndTime

		_, err = client.ScheduledMaintenance().Patch(maintenance.ID, maintenance)
		if err != nil {
			panic(err)
		}

		returnedIncident, err := client.ScheduledMaintenance().CreateStatusUpdate(maintenance.ID, statusUpdate)
		if err != nil {
			panic(err)
		}

		renderedResult, err := renderMaintenance(returnedIncident)
		if err != nil {
			panic(err)
		}

		log.Println(renderedResult)
	},
}
