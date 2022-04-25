package models

import (
	"time"
)

//Region holds information about regions
type Region struct {
	ID          int                    `json:"id"`
	ServiceID   int                    `json:"serviceID"`
	ServiceName string                 `json:"serviceName"`
	Name        string                 `json:"name"`
	RegionCode  string                 `json:"regionCode"`
	Status      ServiceAndRegionStatus `json:"status"`
	Description string                 `json:"description"`
	Tags        []string               `json:"tags"`
	Enabled     bool                   `json:"enabled"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}
