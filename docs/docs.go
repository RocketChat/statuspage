// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Cloud Team",
            "url": "https://open.rocket.chat/group/cloud"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/incidents": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "incident"
                ],
                "summary": "Creates a new incident",
                "operationId": "incident-create",
                "parameters": [
                    {
                        "description": "Incident object",
                        "name": "region",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Incident"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Incident"
                        }
                    }
                }
            }
        },
        "/v1/incidents/{id}/updates": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "incident"
                ],
                "summary": "Creates a new incident update",
                "operationId": "incident-create-update",
                "parameters": [
                    {
                        "description": "Incident update object",
                        "name": "region",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.IncidentUpdate"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "Incident id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.IncidentUpdate"
                        }
                    }
                }
            }
        },
        "/v1/regions": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "region"
                ],
                "summary": "Creates a new region",
                "operationId": "region-create",
                "parameters": [
                    {
                        "description": "Region object",
                        "name": "region",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Region"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Region"
                        }
                    }
                }
            }
        },
        "/v1/regions/{id}": {
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "region"
                ],
                "summary": "Deletes a given region",
                "operationId": "region-delete",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Region id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Incident": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "isMaintenance": {
                    "type": "boolean"
                },
                "latestTweetId": {
                    "type": "integer"
                },
                "maintenance": {
                    "$ref": "#/definitions/models.IncidentMaintenance"
                },
                "originalTweetId": {
                    "type": "integer"
                },
                "regions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.RegionUpdate"
                    }
                },
                "services": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.ServiceUpdate"
                    }
                },
                "status": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "updates": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.IncidentUpdate"
                    }
                }
            }
        },
        "models.IncidentMaintenance": {
            "type": "object",
            "properties": {
                "end": {
                    "type": "string"
                },
                "start": {
                    "type": "string"
                }
            }
        },
        "models.IncidentUpdate": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                },
                "regions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.RegionUpdate"
                    }
                },
                "services": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.ServiceUpdate"
                    }
                },
                "status": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                }
            }
        },
        "models.Region": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "enabled": {
                    "type": "boolean"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "regionCode": {
                    "type": "string"
                },
                "serviceID": {
                    "type": "integer"
                },
                "serviceName": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "models.RegionUpdate": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "regionCode": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "models.ServiceUpdate": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.1",
	Host:        "status.rocket.chat",
	BasePath:    "/api",
	Schemes:     []string{},
	Title:       "Status Central",
	Description: "Operational status for Rocket Chat SaaS & Cloud",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
