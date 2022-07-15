package models

import (
	"time"
)

//ScheduledMaintenance holds the information about the maintenance
type ScheduledMaintenance struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Services    []ServiceUpdate `json:"services,omitempty"`
	Updates     []*StatusUpdate `json:"updates"`

	OriginalTweetID int64 `json:"originalTweetId"`
	LatestTweetID   int64 `json:"latestTweetId"`

	Completed bool `json:"completed"`

	PlannedStart time.Time `json:"plannedStart"`
	PlannedEnd   time.Time `json:"plannedEnd"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
