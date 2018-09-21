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

//IncidentDefaultStatus is the default status of an incident
const IncidentDefaultStatus = "Investigating"
