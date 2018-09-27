package models

import (
	"time"
)

//IncidentUpdate holds an update for the incident
type IncidentUpdate struct {
	ID      int       `json:"id"`
	Time    time.Time `json:"time"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
}
