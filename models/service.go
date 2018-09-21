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

const ServiceStatusOperational = "Operational"
const ServiceStatusDegraded = "Degraded"
const ServiceStatusPartiallyDegraded = "Partially-degraded"
const ServiceStatusOffline = "Offline"
