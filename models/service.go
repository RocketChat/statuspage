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
