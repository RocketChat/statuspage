package core

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/RocketChat/statuscentral/config"
	"github.com/RocketChat/statuscentral/models"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

//GetIncidents retrieves the incidents from the storage layer
func GetIncidents(latest bool) ([]*models.Incident, error) {
	return _dataStore.GetIncidents(latest)
}

//GetIncidentByID retrieves the incident by id, both incident and error will be nil if none found
func GetIncidentByID(id int) (*models.Incident, error) {
	return _dataStore.GetIncidentByID(id)
}

//SendMaintenanceTwitter sends the info about the scheduled maintenance to the offical Rocket.Chat Cloud twitter account.
func SendMaintenanceTwitter(incident *models.Incident) (int64, error) {
	conf := oauth1.NewConfig(config.Config.Twitter.ConsumerKey, config.Config.Twitter.ConsumerSecret)
	token := oauth1.NewToken(config.Config.Twitter.AccessToken, config.Config.Twitter.AccessSecret)
	http := conf.Client(oauth1.NoContext, token)
	http.Timeout = 5 * time.Second

	client := twitter.NewClient(http)
	tmpl, err := template.ParseFiles("templates/incident/tweet/maintenance.tmpl")
	if err != nil {
		return 0, err
	}

	b := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(b, tmpl.Name(), incident)
	if err != nil {
		return 0, err
	}

	tweet, _, err := client.Statuses.Update(b.String(), nil)
	if err != nil {
		return 0, err
	}

	return tweet.ID, nil
}

//SendIncidentTwitter sends the incident info to the offical Rocket.Chat Cloud twitter account.
func SendIncidentTwitter(incident *models.Incident) (int64, error) {
	conf := oauth1.NewConfig(config.Config.Twitter.ConsumerKey, config.Config.Twitter.ConsumerSecret)
	token := oauth1.NewToken(config.Config.Twitter.AccessToken, config.Config.Twitter.AccessSecret)
	http := conf.Client(oauth1.NoContext, token)
	http.Timeout = 5 * time.Second

	client := twitter.NewClient(http)
	tmpl, err := template.ParseFiles("templates/incident/tweet/create.tmpl")
	if err != nil {
		return 0, err
	}

	b := &bytes.Buffer{}
	if err := tmpl.ExecuteTemplate(b, tmpl.Name(), incident); err != nil {
		return 0, err
	}

	tweet, _, err := client.Statuses.Update(b.String(), nil)
	if err != nil {
		return 0, err
	}

	return tweet.ID, nil
}

//SendIncidentUpdateTwitter sends the incident update info to the offical Rocket.Chat Cloud twitter account.
func SendIncidentUpdateTwitter(incident *models.Incident, update *models.IncidentUpdate) (int64, error) {
	conf := oauth1.NewConfig(config.Config.Twitter.ConsumerKey, config.Config.Twitter.ConsumerSecret)
	token := oauth1.NewToken(config.Config.Twitter.AccessToken, config.Config.Twitter.AccessSecret)
	http := conf.Client(oauth1.NoContext, token)
	http.Timeout = 5 * time.Second

	client := twitter.NewClient(http)
	tmpl, err := template.ParseFiles("templates/incident/tweet/update.tmpl")
	if err != nil {
		return 0, err
	}

	b := &bytes.Buffer{}
	if err := tmpl.ExecuteTemplate(b, tmpl.Name(), update); err != nil {
		return 0, err
	}

	var params = &twitter.StatusUpdateParams{}
	if incident.OriginalTweetID != 0 {
		params.InReplyToStatusID = incident.OriginalTweetID
	}

	tweet, _, err := client.Statuses.Update(fmt.Sprintf("%s", b.String()), params)
	if err != nil {
		return 0, err
	}

	return tweet.ID, nil
}

//CreateIncident creates the incident in the storage layer
func CreateIncident(incident *models.Incident) (*models.Incident, error) {
	ensureIncidentDefaults(incident)

	if incident.Status == models.IncidentStatusScheduledMaintenance {
		incident.IsMaintenance = true
	}

	if len(incident.Updates) == 0 {
		if incident.IsMaintenance {
			update := models.IncidentUpdate{
				Time:   incident.Time,
				Status: incident.Status,
				Message: fmt.Sprintf("Starts at %s with a scheduled end at %s",
					incident.Maintenance.Start.Format(time.RFC1123Z),
					incident.Maintenance.End.Format(time.RFC1123Z)),
			}

			incident.Updates = append(incident.Updates, &update)
		} else {
			update := models.IncidentUpdate{
				Time:    incident.Time,
				Status:  incident.Status,
				Message: "Initial status of " + incident.Status.String(),
			}

			incident.Updates = append(incident.Updates, &update)
		}
	}

	if incident.IsMaintenance {
		for _, s := range incident.Services {
			if err := updateServiceToStatus(s.Name, models.ServiceStatusScheduledMaintenance); err != nil {
				return nil, err
			}
		}
	} else {
		for _, s := range incident.Services {
			if err := updateServiceToStatus(s.Name, s.Status); err != nil {
				return nil, err
			}
		}
	}

	if err := _dataStore.CreateIncident(incident); err != nil {
		return nil, err
	}

	if config.Config.Twitter.Enabled {
		if incident.Status == models.IncidentStatusScheduledMaintenance {
			tweetID, err := SendMaintenanceTwitter(incident)
			if err == nil {
				incident.OriginalTweetID = tweetID
				if err := _dataStore.UpdateIncident(incident); err != nil {
					return nil, err
				}
			}
		} else {
			tweetID, err := SendIncidentTwitter(incident)
			if err == nil {
				incident.OriginalTweetID = tweetID
				if err := _dataStore.UpdateIncident(incident); err != nil {
					return nil, err
				}
			}
		}
	}

	return incident, nil
}

//DeleteIncident removes the incident from the storage layer
func DeleteIncident(id int) error {
	return _dataStore.DeleteIncident(id)
}

//CreateIncidentUpdate creates an update for an incident
func CreateIncidentUpdate(incidentID int, update *models.IncidentUpdate) (*models.Incident, error) {
	if incidentID <= 0 {
		return nil, errors.New("invalid incident id")
	}

	if update.Message == "" {
		return nil, errors.New("message property is missing")
	}

	if update.Status == "" {
		return nil, errors.New("status property is missing")
	}

	status, ok := models.IncidentStatuses[strings.ToLower(update.Status.String())]
	if !ok {
		return nil, errors.New("invalid status value")
	}

	update.Status = status

	if err := _dataStore.CreateIncidentUpdate(incidentID, update); err != nil {
		return nil, err
	}

	incident, err := _dataStore.GetIncidentByID(incidentID)
	if err != nil {
		return nil, err
	}

	if status != models.IncidentStatusResolved {
		for _, s := range update.Services {
			if err := updateServiceToStatus(s.Name, s.Status); err != nil {
				return nil, err
			}
		}
	} else {
		for _, s := range incident.Services {
			if err := updateServiceToStatus(s.Name, models.ServiceStatusNominal); err != nil {
				return nil, err
			}
		}
	}

	if config.Config.Twitter.Enabled {
		tweetID, err := SendIncidentUpdateTwitter(incident, update)
		if err == nil && tweetID != 0 {
			incident.LatestTweetID = tweetID
			if err := _dataStore.UpdateIncident(incident); err != nil {
				return nil, err
			}
		}
	}

	return incident, nil
}

func ensureIncidentDefaults(incident *models.Incident) {
	if incident.Updates == nil {
		incident.Updates = make([]*models.IncidentUpdate, 0)
	}

	if incident.Status == "" {
		incident.Status = models.IncidentDefaultStatus
	}

	if incident.Title == "" {
		incident.Title = "Unknown"
	}

	if incident.Time.IsZero() {
		incident.Time = time.Now()
	}
}

//AggregateIncidents aggregates the incidents
func AggregateIncidents(incidents []*models.Incident) models.AggregatedIncidents {
	aggregatedIncidents := models.AggregatedIncidents{}

	for i := 0; i < config.Config.Website.DaysToAggregate; i++ {
		t := time.Now().Add(-time.Duration(i) * 24 * time.Hour)
		filteredIncidents := []*models.Incident{}

		for _, incident := range incidents {
			if incident.Time.Day() == t.Day() {
				filteredIncidents = append(filteredIncidents, incident)
			}
		}

		aggregatedIncidents = append(aggregatedIncidents, models.AggregatedIncident{
			Time:      t,
			Incidents: filteredIncidents,
		})
	}

	return aggregatedIncidents
}
