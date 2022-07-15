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

//AggregatedScheduledMaintenance contains the scheduled maintenance aggregated
type AggregatedScheduledMaintenance struct {
	Time                 time.Time
	ScheduledMaintenance []*ScheduledMaintenance
}

//AggregatedScheduledMaintenances holds several aggregated scheduled maintenances
type AggregatedScheduledMaintenances []AggregatedScheduledMaintenance
