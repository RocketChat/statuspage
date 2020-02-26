package models

//ServiceUpdate is for a service update
type ServiceUpdate struct {
	Name   string        `json:"name"`
	Status ServiceStatus `json:"status"`
}
