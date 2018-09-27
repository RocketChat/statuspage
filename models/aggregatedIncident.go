package models

import (
	"time"
)

//AggregatedIncident contains the incidents aggregated
type AggregatedIncident struct {
	Time      time.Time
	Incidents []*Incident
}

//AggregatedIncidents holds several aggregated incidents
type AggregatedIncidents []AggregatedIncident
