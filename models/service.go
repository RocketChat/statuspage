package models

//Service holds information about the service
type Service struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Status      string   `json:"status"`
	Description string   `json:"description"`
	Group       string   `json:"group"`
	Link        string   `json:"link"`
	Tags        []string `json:"tags"`
	Enabled     bool     `json:"enabled"`
}
