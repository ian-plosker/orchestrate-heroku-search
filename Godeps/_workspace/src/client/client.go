// Package client implements a basic Orchestrate client
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	r "reflect"
)

const (
	rootUri = "https://api.orchestrate.io/v0/"
)

type Client struct {
	HttpClient *http.Client
	AuthToken  string
}

type OrchestrateError struct {
	Status  string `json:"-"`
	StatusCode int `json:"-"`
	Message string `json:"message"`
	Code string 	 `json:"code"`
}

// NewClient returns a new orchestrate client.
func NewClient(authToken string) *Client {
	httpClient := &http.Client{}

	return &Client{
		HttpClient: httpClient,
		AuthToken:  authToken,
	}
}

func newError(resp *http.Response) error {
	decoder := json.NewDecoder(resp.Body)
	orchestrateError := new(OrchestrateError)
	decoder.Decode(orchestrateError)

	orchestrateError.Status = resp.Status
	orchestrateError.StatusCode = resp.StatusCode

	return orchestrateError
}

func (e *OrchestrateError) Error() string {
	return fmt.Sprintf(`%v: %v`, e.Status, e.Message)
}

func (client Client) doRequest(method, trailingPath string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, rootUri+trailingPath, body)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(client.AuthToken, "")

	if method == "PUT" {
		req.Header.Add("Content-Type", "application/json")
	}

	return client.HttpClient.Do(req)
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
