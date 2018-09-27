package core

import (
	"errors"
	"strings"
	"time"

	"github.com/RocketChat/statuscentral/config"
	"github.com/RocketChat/statuscentral/models"
)

//GetIncidents retrieves the incidents from the storage layer
func GetIncidents(latest bool) ([]*models.Incident, error) {
	return _dataStore.GetIncidents(latest)
}

//GetIncidentByID retrieves the incident by id, both incident and error will be nil if none found
func GetIncidentByID(id int) (*models.Incident, error) {
	return _dataStore.GetIncidentByID(id)
}

//CreateIncident creates the incident in the storage layer
func CreateIncident(incident *models.Incident) error {
	ensureIncidentDefaults(incident)

	if len(incident.Updates) == 0 {
		update := models.IncidentUpdate{
			Time:    time.Now(),
			Status:  incident.Status,
			Message: "Initial status of " + incident.Status,
		}

		incident.Updates = append(incident.Updates, &update)
	}

	return _dataStore.CreateIncident(incident)
}

//UpdateIncident updates the incident in the storage layer
func UpdateIncident(incident *models.Incident) error {
	ensureIncidentDefaults(incident)

	return _dataStore.UpdateIncident(incident)
}

//DeleteIncident removes the incident from the storage layer
func DeleteIncident(id int) error {
	return _dataStore.DeleteIncident(id)
}

//CreateIncidentUpdate creates an update for an incident
func CreateIncidentUpdate(incidentID int, update *models.IncidentUpdate) error {
	if incidentID <= 0 {
		return errors.New("invalid incident id")
	}

	if update.Message == "" {
		return errors.New("message property is missing")
	}

	if update.Status == "" {
		return errors.New("status property is missing")
	}

	status, ok := models.IncidentStatuses[strings.ToLower(update.Status)]
	if !ok {
		return errors.New("invalid status value")
	}

	update.Status = status

	return _dataStore.CreateIncidentUpdate(incidentID, update)
}

func ensureIncidentDefaults(incident *models.Incident) {
	if incident.Updates == nil {
		incident.Updates = make([]*models.IncidentUpdate, 0)
	}

	if incident.Status == "" {
		incident.Status = models.IncidentDefaultStatus
	}

	if incident.Title == "" {
		incident.Title = "Unknown"
	}

	if incident.Time.IsZero() {
		incident.Time = time.Now()
	}
}

//AggregateIncidents aggregates the incidents
func AggregateIncidents(incidents []*models.Incident) models.AggregatedIncidents {
	aggregatedIncidents := models.AggregatedIncidents{}

	for i := 0; i < config.Config.Website.DaysToAggregate; i++ {
		t := time.Now().Add(-time.Duration(i) * 24 * time.Hour)
		filteredIncidents := []*models.Incident{}

		for _, incident := range incidents {
			if incident.Time.Day() == t.Day() {
				filteredIncidents = append(filteredIncidents, incident)
			}
		}

		aggregatedIncidents = append(aggregatedIncidents, models.AggregatedIncident{
			Time:      t,
			Incidents: filteredIncidents,
		})
	}

	return aggregatedIncidents
}
