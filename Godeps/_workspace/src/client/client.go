// Copyright 2014, Orchestrate.IO, Inc.

// A client for use with Orchestrate.io: http://orchestrate.io/
//
// Orchestrate unifies multiple databases through one simple REST API.
// Orchestrate runs as a service and supports queries like full-text
// search, events, graph, and key/value.
//
// You can sign up for an Orchestrate account here:
// http://dashboard.orchestrate.io
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	r "reflect"
	"time"
)

// The root path for all API endpoints.
const rootUri = "https://api.orchestrate.io/v0/"

var (
	Timeout = 5 * time.Second
	Transport http.RoundTripper = &http.Transport{
		MaxIdleConnsPerHost: 100,
		ResponseHeaderTimeout: Timeout,
		Dial: func (network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, Timeout)
		},
	}
)

type Client struct {
	httpClient *http.Client
	authToken  string
}

// An implementation of 'error' that exposes all the orchestrate specific
// error details.
type OrchestrateError struct {
	Status  string `json:"-"`
	StatusCode int `json:"-"`
	Message string `json:"message"`
	Code string    `json:"code"`
}

// Returns a new Client object that will use the given authToken for
// authorization against Orchestrate. This token can be obtained
// at http://dashboard.orchestrate.io
func NewClient(authToken string) *Client {
	return &Client{
		httpClient: &http.Client{Transport: Transport},
		authToken:  authToken,
	}
}

// Creates a new OrchestrateError from a given http.Response object.
func newError(resp *http.Response) error {
	decoder := json.NewDecoder(resp.Body)
	orchestrateError := new(OrchestrateError)
	decoder.Decode(orchestrateError)

	orchestrateError.Status = resp.Status
	orchestrateError.StatusCode = resp.StatusCode

	return orchestrateError
}

func (e OrchestrateError) Error() string {
	return fmt.Sprintf(`%v (%v): %v`, e.Status, e.StatusCode, e.Message)
}

func (client *Client) doRequest(method, trailingPath string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, rootUri+trailingPath, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(client.authToken, "")

	if method == "PUT" {
		req.Header.Add("Content-Type", "application/json")
	}

	return client.httpClient.Do(req)
}

func ValueToStruct(value map[string]interface{}, dest interface{}) bool {
	structVal := r.Indirect(r.ValueOf(dest))
	structType := structVal.Type()

	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		name := structField.Name

		if jField := structField.Tag.Get("json"); jField != "" {
			name = jField
		}

		if fieldValue, present := value[name]; present {
			fieldVal := r.ValueOf(fieldValue)

			if fieldVal.Type() != structField.Type {
				fieldVal = fieldVal.Convert(structField.Type)
			}

			structVal.FieldByName(structField.Name).Set(fieldVal)
		}
	}

	return true
}
