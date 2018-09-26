package core

import (
	"time"

	"github.com/RocketChat/statuspage/config"
	"github.com/RocketChat/statuspage/models"
)

//GetIncidents retrieves the incidents from the storage layer
func GetIncidents(latest bool) ([]*models.Incident, error) {
	return _dataStore.GetIncidents(latest)
}

//CreateIncident creates the incident in the storage layer
func CreateIncident(incident *models.Incident) error {
	ensureIncidentDefaults(incident)

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
