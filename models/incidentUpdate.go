package models

import (
	"time"
)

//StatusUpdate holds an update for an incident or scheduled maintenance
type StatusUpdate struct {
	ID       int             `json:"id"`
	Time     time.Time       `json:"time"`
	Status   IncidentStatus  `json:"status"`
	Message  string          `json:"message"`
	Services []ServiceUpdate `json:"services,omitempty"`
	Regions  []RegionUpdate  `json:"regions,omitempty"`
}
