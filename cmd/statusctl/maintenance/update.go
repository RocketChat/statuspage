package maintenance

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/RocketChat/statuscentral/cmd/statusctl/common"
	"github.com/RocketChat/statuscentral/models"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "add maintenance update",
	Example: "statusctl maintenance update",
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

func updateImpactedServices(services []models.ServiceUpdate) ([]models.ServiceUpdate, error) {
	updatingServices := true
	servicesUpdated := map[string]bool{}

	for updatingServices {
		for i, service := range services {
			serviceLine := fmt.Sprintf("%d) %s - %s", i, service.Name, service.Status)

			if servicesUpdated[service.Name] {
				serviceLine = serviceLine + " ** not yet submitted"
			}

			log.Println(serviceLine)
		}

		service, err := common.IntPrompt("Select a service to Update [0]:", 0)
		if err != nil {
			return services, errors.New("invalid selection")
		}

		currentServiceStatus := 1

		for j, serviceStatus := range models.ServiceStatusArray {
			if serviceStatus == services[service].Status {
				currentServiceStatus = j
			}

			log.Printf("%d) %s\n", j, serviceStatus)
		}

		serviceStatus, err := common.IntPrompt(fmt.Sprintf("Select service status [%d]:", currentServiceStatus), currentServiceStatus)
		if err != nil {
			return services, errors.New("invalid selection")
		}

		services[service] = models.ServiceUpdate{
			Name:    services[service].Name,
			Status:  models.ServiceStatusArray[serviceStatus],
			Regions: []string{},
		}

		servicesUpdated[services[service].Name] = true

		more, err := common.GetYesNoPrompt("Update another service?", true)
		if err != nil {
			return services, err
		}

		if !more {
			updatingServices = false
		}
	}

	log.Println(services)

	return services, nil
}
