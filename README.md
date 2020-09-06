# statuscentral
> Self hosted status page written in golang!

## Creating an Incident
First step is to create an incident and describe which services were affected and what those services status is.

### Service Statuses
* Nominal
* Degraded
* Partial-outage
* Outage
* Scheduled Maintenance
* Unknown

### Incident Status
* Investigating
* Identified
* Update
* Monitoring
* Resolved

### Incident Creation Call
`POST https://status.rocket.chat/api/v1/incidents`

Request body:
```json
{
	"title": "Slowness Reported Again",
    "time": "2020-01-22T14:39:24.495623-06:00",
    "status": "Investigating",
	"services": [
		{
			"name": "Marketplace",
			"status": "Degraded"
		}
	]
}
```

Resulting object:
```json
{
  "id": 2,
  "time": "2020-01-22T14:39:24.495623-06:00",
  "title": "Slowness Reported Again",
  "status": "Investigating",
  "updates": [
    {
      "id": 0,
      "time": "2020-02-25T18:44:35.592427-06:00",
      "status": "Investigating",
      "message": "Initial status of Investigating"
    }
  ],
  "updatedAt": "2020-02-25T18:44:35.604079-06:00"
}
```

### Incident Update
`POST https://status.rocket.chat/api/v1/incidents/:id/updates`
```json
{
	"message": "Testing msg",
	"status": "Identified",
	"time": "2020-02-25T19:00:22.585515764-05:00",
    "serivces": [
        {
            "name": "Marketplace",
            "status": "Partial-outage"
        }
    ]
}
```