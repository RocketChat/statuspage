package models

import (
	"time"
)

//Incident holds the information about the incident
type Incident struct {
	ID          int                 `json:"id"`
	Time        time.Time           `json:"time"`
	Title       string              `json:"title"`
	Status      IncidentStatus      `json:"status"`
	Services    []ServiceUpdate     `json:"services,omitempty"`
	Updates     []*IncidentUpdate   `json:"updates"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Maintenance IncidentMaintenance `json:"maintenance"`
}

//IncidentMaintenance contains the data about a scheduled maintenance.
type IncidentMaintenance struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

//IncidentStatus represents the status of the incident
type IncidentStatus string

func (is IncidentStatus) String() string {
	return string(is)
}

const (
	//IncidentStatusScheduledMaintenance - Schedule maintenance Incident
	IncidentStatusScheduledMaintenance IncidentStatus = "Scheduled Maintenance"
	//IncidentStatusInvestigating - Investigating Incident
	IncidentStatusInvestigating IncidentStatus = "Investigating"
	//IncidentStatusIdentified - Identified cause of Incident
	IncidentStatusIdentified IncidentStatus = "Identified"
	//IncidentStatusUpdate - An update to the Incident. Does not update overall incident status
	IncidentStatusUpdate IncidentStatus = "Update"
	//IncidentStatusMonitoring - Monitoring typically after a fix applied
	IncidentStatusMonitoring IncidentStatus = "Monitoring"
	//IncidentStatusResolved - Resolved the incident
	IncidentStatusResolved IncidentStatus = "Resolved"
	//IncidentDefaultStatus is the default status of an incident
	IncidentDefaultStatus IncidentStatus = IncidentStatusInvestigating
)

//IncidentStatuses holds all of the valid incident statuses
var IncidentStatuses = map[string]IncidentStatus{
	"scheduled maintenance": IncidentStatusScheduledMaintenance,
	"investigating":         IncidentStatusInvestigating,
	"identified":            IncidentStatusIdentified,
	"update":                IncidentStatusUpdate,
	"monitoring":            IncidentStatusMonitoring,
	"resolved":              IncidentStatusResolved,
	"default":               IncidentDefaultStatus,
}
