// Copyright 2014, Orchestrate.IO, Inc.

// A client for use with Orchestrate.io: http://orchestrate.io/
//
// Orchestrate unifies multiple databases through one simple REST API.
// Orchestrate runs as a service and supports queries like full-text
// search, events, graph, and key/value.
//
// You can sign up for an Orchestrate account here:
// http://dashboard.orchestrate.io
package gorc

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	// The root path for all API endpoints.
	rootUri = "https://api.orchestrate.io/v0/"
)

var (
	// The default timeout that will be used for connections. This is used
	// with the default Transport to establish how long a connection attempt
	// can take. This is not the data transfer timeout. Changing this will
	// impact all new connections made with the default transport.
	DefaultDialTimeout = 3 * time.Second

	// This is the default http.Transport that will be associated with new
	// clients. If overwritten then only new clients will be impacted, old
	// clients will continue to use the pre-existing transport.
	DefaultTransport *http.Transport = &http.Transport{
		// In the default configuration we allow 4 idle connections to the
		// api server. This limits the number of live connections to our
		// load balancer which reduces load. If needed this can be increased
		// for high volume clients.
		MaxIdleConnsPerHost: 4,

		// This timeout value is how long the http client library will wait
		// for data before abandoning the call. If this is set too low then
		// high work calls, or high latency connections can trip timeouts
		// too often.
		ResponseHeaderTimeout: 3 * time.Second,

		// The default Dial function is over written so it uses net.DialTimeout
		// instead.
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, DefaultDialTimeout)
		},
	}
)

// An Orchestrate Client object.
type Client struct {
	httpClient *http.Client
	authToken  string
}

// An implementation of 'error' that exposes all the orchestrate specific
// error details.
type OrchestrateError struct {
	// The status string returned from the HTTP call.
	Status string `json:"-"`

	// The status, as an integer, returned from the HTTP call.
	StatusCode int `json:"-"`

	// The Orchestrate specific message representing the error.
	Message string `json:"message"`
}

// A representation of a Key/Value object's path within Orchestrate.
type Path struct {
	Collection string `json:"collection"`
	Key        string `json:"key"`
	Ref        string `json:"ref"`
}

// Returns a new Client object that will use the given authToken for
// authorization against Orchestrate. This token can be obtained
// at http://dashboard.orchestrate.io
func NewClient(authToken string) *Client {
	return NewClientWithTransport(authToken, DefaultTransport)
}

// Like NewClient, except that it allows a specific http.Transport to be
// provided for use, rather than DefaultTransport.
func NewClientWithTransport(authToken string, transport *http.Transport) *Client {
	return &Client{
		httpClient: &http.Client{Transport: transport},
		authToken:  authToken,
	}
}

// Creates a new OrchestrateError from a given http.Response object.
func newError(resp *http.Response) error {
	oe := &OrchestrateError{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
	}
	if data, err := ioutil.ReadAll(resp.Body); err != nil {
		oe.Message = fmt.Sprintf("Can not read HTTP response: %s", err)
		return oe
	} else if err := json.Unmarshal(data, oe); err != nil {
		oe.Message = fmt.Sprintf("Can not unmarshal JSON response '''%s''': %s", string(data), err)
		return oe
	}

	return oe
}

func (e OrchestrateError) Error() string {
	return fmt.Sprintf("%s (%d): %s", e.Status, e.StatusCode, e.Message)
}

// Executes an HTTP request.
func (c *Client) doRequest(method, trailing string, headers map[string]string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, rootUri+trailing, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.authToken, "")

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if method == "PUT" {
		req.Header.Add("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}
