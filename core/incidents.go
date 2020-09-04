package core

import (
	"errors"
	"fmt"
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

	if incident.Status == models.IncidentStatusScheduledMaintenance {
		incident.IsMaintenance = true
	}

	if len(incident.Updates) == 0 {
		if incident.IsMaintenance {
			update := models.IncidentUpdate{
				Time:   incident.Time,
				Status: incident.Status,
				Message: fmt.Sprintf("Starts at %s with a scheduled end at %s",
					incident.Maintenance.Start.Format(time.RFC1123Z),
					incident.Maintenance.End.Format(time.RFC1123Z)),
			}

			incident.Updates = append(incident.Updates, &update)
		} else {
			update := models.IncidentUpdate{
				Time:    incident.Time,
				Status:  incident.Status,
				Message: "Initial status of " + incident.Status.String(),
			}

			incident.Updates = append(incident.Updates, &update)
		}
	}

	if incident.IsMaintenance {
		for _, s := range incident.Services {
			if err := updateServiceToStatus(s.Name, models.ServiceStatusScheduledMaintenance); err != nil {
				return err
			}
		}
	} else {
		for _, s := range incident.Services {
			if err := updateServiceToStatus(s.Name, s.Status); err != nil {
				return err
			}
		}
	}

	return _dataStore.CreateIncident(incident)
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

	status, ok := models.IncidentStatuses[strings.ToLower(update.Status.String())]
	if !ok {
		return errors.New("invalid status value")
	}

	update.Status = status

	if err := _dataStore.CreateIncidentUpdate(incidentID, update); err != nil {
		return err
	}

	if status != models.IncidentStatusResolved {
		for _, s := range update.Services {
			if err := updateServiceToStatus(s.Name, s.Status); err != nil {
				return err
			}
		}

		return nil
	}

	incident, err := _dataStore.GetIncidentByID(incidentID)
	if err != nil {
		return err
	}

	for _, s := range incident.Services {
		if err := updateServiceToStatus(s.Name, models.ServiceStatusNominal); err != nil {
			return err
		}
	}

	return nil
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
