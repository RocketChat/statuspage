package models

import (
	"time"
)

//Service holds information about the service
type Service struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Group       string    `json:"group"`
	Link        string    `json:"link"`
	Tags        []string  `json:"tags"`
	Enabled     bool      `json:"enabled"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

//ServiceStatusOperational - Everything is good
const ServiceStatusOperational = "Operational"

//ServiceStatusDegraded - Degraded Performance
const ServiceStatusDegraded = "Degraded"

//ServiceStatusPartialOutage - Partial Outage
const ServiceStatusPartialOutage = "Partial-outage"

//ServiceStatusOutage - Outage
const ServiceStatusOutage = "Outage"

//ServiceStatusScheduledMaintenance - Scheduled Maintenance
const ServiceStatusScheduledMaintenance = "Scheduled Maintenance"

//ServiceStatusUnknown - Unknown - used when the services were loaded from config
const ServiceStatusUnknown = "Unknown"

//ServiceStatusValues holds the names to numbers for values of the status
var ServiceStatusValues = map[string]int{
	"Operational":           0,
	"Degraded":              1,
	"Partial-outage":        2,
	"Outage":                3,
	"Scheduled Maintenance": 4,
	"Unknown":               5,
}
