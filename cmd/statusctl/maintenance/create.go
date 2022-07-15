package maintenance

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/RocketChat/statuscentral/cmd/statusctl/common"
	"github.com/RocketChat/statuscentral/models"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "create incident",
	Example: "statusctl maintenance create",
	Run: func(c *cobra.Command, args []string) {
		client := common.GetStatusCentralClient()

		services, err := client.Services().GetMultiple()
		if err != nil {
			panic(err)
		}

		log.Println("Great!  Lets get that scheduled!")

		title := common.StringPrompt("Maintenance Short Descripton / Title:")
		description := common.StringPrompt("Longer Description:")

		servicesImpacted, err := getImpactedServices(services)
		if err != nil {
			panic(err)
		}

		plannedStartTimeText := common.StringPrompt("Planned Start UTC Time (format: 2022/07/15 17:08:12):")

		plannedStartTime, err := time.ParseInLocation("2006/01/02 15:04:05", plannedStartTimeText, time.UTC)
		if err != nil {
			panic(err)
		}

		plannedEndTimeText := common.StringPrompt("Planned End UTC Time (format: 2022/07/15 17:08:12):")

		plannedEndTime, err := time.ParseInLocation("2006/01/02 15:04:05", plannedEndTimeText, time.UTC)
		if err != nil {
			panic(err)
		}

		scheduledMaintenance := &models.ScheduledMaintenance{
			Title:        title,
			Description:  description,
			Services:     servicesImpacted,
			PlannedStart: plannedStartTime,
			PlannedEnd:   plannedEndTime,
		}

		returnedScheduledMaintenance, err := client.ScheduledMaintenance().Create(scheduledMaintenance)
		if err != nil {
			panic(err)
		}

		log.Println(fmt.Sprintf("Maintenance %d Scheduled!", returnedScheduledMaintenance.ID))

		rendered, err := renderMaintenance(returnedScheduledMaintenance)
		if err != nil {
			panic(err)
		}

		log.Println(rendered)
	},
}

func getImpactedServices(services []*models.Service) ([]models.ServiceUpdate, error) {
	gettingServices := true
	serviceUpdates := []models.ServiceUpdate{}

	for gettingServices {
		for i, service := range services {
			log.Printf("%d) %s\n", i, service.Name)
		}

		service, err := common.IntPrompt("Select a service impacted [1]:", 1)
		if err != nil {
			return serviceUpdates, errors.New("invalid selection")
		}

		for i, serviceStatus := range models.ServiceStatusArray {
			log.Printf("%d) %s\n", i, serviceStatus)
		}

		serviceStatus, err := common.IntPrompt("Select service status [1]:", 1)
		if err != nil {
			return serviceUpdates, errors.New("invalid selection")
		}

		serviceUpdates = append(serviceUpdates, models.ServiceUpdate{
			Name:    services[service].Name,
			Status:  models.ServiceStatusArray[serviceStatus],
			Regions: []string{},
		})

		more, err := common.GetYesNoPrompt("Add another service?", true)
		if err != nil {
			return serviceUpdates, err
		}

		if !more {
			gettingServices = false
		}
	}

	return serviceUpdates, nil
}
