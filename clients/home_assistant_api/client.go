package home_assistant_api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"smarthome_ai_bot/clients"
	"strings"
)

type clientImpl struct {
	apiHost     string
	bearerToken string
}

type Config struct {
	ApiHost     string
	BearerToken string
}

func NewClient(cfg *Config) (clients.SmarthomeAPI, error) {
	if cfg == nil {
		return nil, errors.New("missing parameter: cfg")
	}

	if cfg.ApiHost == "" {
		return nil, errors.New("missing parameter: cfg.ApiHost")
	}

	if cfg.BearerToken == "" {
		return nil, errors.New("missing parameter: cfg.BearerToken")
	}

	return &clientImpl{
		apiHost:     cfg.ApiHost,
		bearerToken: cfg.BearerToken,
	}, nil
}

type performServiceRequest struct {
	EntityID      string `json:"entity_id"`
	Value         string `json:"value,omitempty"`
	BrightnessPCT string `json:"brightness_pct,omitempty"`
}

func (client *clientImpl) PerformService(ctx context.Context, service, entityID, setValue string) (string, error) {
	// replace all dots in service with a forward slash
	service = strings.ReplaceAll(service, ".", "/")

	serviceUrl := "/api/services/" + service

	url := client.apiHost + serviceUrl

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + client.bearerToken

	// make an io reader for the json body
	var postBody performServiceRequest

	postBody = performServiceRequest{
		EntityID: entityID,
	}

	switch service {
	case "light/turn_on":
		postBody.BrightnessPCT = setValue
	case "input_number/set_value":
		postBody.Value = setValue
	default:
		postBody.Value = setValue
	}

	jsonBody, err := json.Marshal(postBody)
	if err != nil {
		return "", err
	}

	bodyReader := bytes.NewReader(jsonBody)

	log.Printf("url: %s json: %s", serviceUrl, jsonBody)

	// Create a new request using http
	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return "", err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "system failed to execute command", nil
	}

	return "system executed command successfully", nil
}

type queryResponse struct {
	State string `json:"state"`
}

func (client *clientImpl) QueryEntity(ctx context.Context, entityID string) (string, error) {
	url := client.apiHost + "/api/states/" + entityID

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + client.bearerToken

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response queryResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.State, nil
}
