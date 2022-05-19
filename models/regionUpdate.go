package models

//RegionUpdate is for a region update
type RegionUpdate struct {
	Name       string                 `json:"name"`
	Status     ServiceAndRegionStatus `json:"status"`
	RegionCode string                 `json:"regionCode"`
}
