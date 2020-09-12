package models

import (
	"strings"
	"time"
)

//Service holds information about the service
type Service struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Status      ServiceStatus `json:"status"`
	Description string        `json:"description"`
	Group       string        `json:"group"`
	Link        string        `json:"link"`
	Tags        []string      `json:"tags"`
	Enabled     bool          `json:"enabled"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

//ServiceStatus represents the status of a service
type ServiceStatus string

func (ss ServiceStatus) String() string {
	return string(ss)
}

//ToLower converts the status to lowercase string
func (ss ServiceStatus) ToLower() string {
	return strings.ToLower(ss.String())
}

const (
	//ServiceStatusNominal - Everything is good
	ServiceStatusNominal ServiceStatus = "Nominal"
	//ServiceStatusDegraded - Degraded Performance
	ServiceStatusDegraded ServiceStatus = "Degraded"
	//ServiceStatusPartialOutage - Partial Outage
	ServiceStatusPartialOutage ServiceStatus = "Partial-outage"
	//ServiceStatusOutage - Outage
	ServiceStatusOutage ServiceStatus = "Outage"
	//ServiceStatusScheduledMaintenance - Scheduled Maintenance
	ServiceStatusScheduledMaintenance ServiceStatus = "Scheduled Maintenance"
	//ServiceStatusUnknown - Unknown - used when the services were loaded from config
	ServiceStatusUnknown ServiceStatus = "Unknown"
)

//ServiceStatusValues holds the names to numbers for values of the status
var ServiceStatusValues = map[string]int{
	"Nominal":               0,
	"Degraded":              1,
	"Partial-outage":        2,
	"Outage":                3,
	"Scheduled Maintenance": 4,
	"Unknown":               5,
}

//ServiceStatuses holds a map of the lower case service statuses
var ServiceStatuses = map[string]ServiceStatus{
	ServiceStatusNominal.ToLower():              ServiceStatusNominal,
	ServiceStatusDegraded.ToLower():             ServiceStatusDegraded,
	ServiceStatusPartialOutage.ToLower():        ServiceStatusPartialOutage,
	ServiceStatusOutage.ToLower():               ServiceStatusOutage,
	ServiceStatusScheduledMaintenance.ToLower(): ServiceStatusScheduledMaintenance,
	ServiceStatusUnknown.ToLower():              ServiceStatusUnknown,
}
