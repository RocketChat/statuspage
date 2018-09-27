package models

import (
	"time"
)

//Incident holds the information about the incident
type Incident struct {
	ID        int               `json:"id"`
	Time      time.Time         `json:"time"`
	Title     string            `json:"title"`
	Status    string            `json:"status"`
	Updates   []*IncidentUpdate `json:"updates"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

//IncidentStatusInvestigating - Investigating Incident
const IncidentStatusInvestigating = "Investigating"

//IncidentStatusIdentified - Identified cause of Incident
const IncidentStatusIdentified = "Identified"

//IncidentStatusUpdate - An update to the Incident. Does not update overall incident status
const IncidentStatusUpdate = "Update"

//IncidentStatusMonitoring - Monitoring typically after a fix applied
const IncidentStatusMonitoring = "Monitoring"

//IncidentStatusResolved - Resolved the incident
const IncidentStatusResolved = "Resolved"

//IncidentDefaultStatus is the default status of an incident
const IncidentDefaultStatus = IncidentStatusInvestigating

//IncidentStatuses holds all of the valid incident statuses
var IncidentStatuses = map[string]string{
	"investigating": "Investigating",
	"identified":    "Identified",
	"update":        "Update",
	"monitoring":    "Monitoring",
	"resolved":      "Resolved",
	"default":       IncidentDefaultStatus,
}
