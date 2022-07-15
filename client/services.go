package client

import (
	"github.com/RocketChat/statuscentral/models"
)

// ServicesInterface services interface
type ServicesInterface interface {
	GetMultiple() (result []*models.Service, err error)
}

type services struct {
	client *Client
}

func (s *services) GetMultiple() (result []*models.Service, err error) {
	req, err := s.client.buildRequest("GET", "/api/v1/services", nil)

	if err != nil {
		return nil, err
	}

	result = []*models.Service{}

	resp, err := s.client.do(req, &result)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}
