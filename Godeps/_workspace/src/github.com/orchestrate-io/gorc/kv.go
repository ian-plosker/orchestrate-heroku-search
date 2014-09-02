// Copyright 2014, Orchestrate.IO, Inc.

package gorc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

// Holds results returned from a KV list query.
type KVResults struct {
	Count   uint64     `json:"count"`
	Results []KVResult `json:"results"`
	Next    string     `json:"next,omitempty"`
}

// An individual Key/Value result.
type KVResult struct {
	Path     Path            `json:"path"`
	RawValue json.RawMessage `json:"value"`
}

// Get a collection-key pair's value.
func (c *Client) Get(collection, key string) (*KVResult, error) {
	return c.GetPath(&Path{Collection: collection, Key: key})
}

// Get the value at a path.
func (c *Client) GetPath(path *Path) (*KVResult, error) {
	resp, err := c.doRequest("GET", path.trailingGetURI(), nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, newError(resp)
	}

	// TODO: Check for a content-length header so we can pre-allocate buffer
	// space.
	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	if path.Ref == "" {
		parts := strings.SplitAfter(resp.Header.Get("Content-Location"), "/")
		if len(parts) >= 6 {
			path.Ref = parts[5]
		}
	}

	return &KVResult{Path: *path, RawValue: buf.Bytes()}, nil
}

// Store a value to a collection-key pair.
func (c *Client) Put(collection string, key string, value interface{}) (*Path, error) {
	reader, writer := io.Pipe()
	encoder := json.NewEncoder(writer)

	go func() { writer.CloseWithError(encoder.Encode(value)) }()
	return c.PutRaw(collection, key, reader)
}

// Store a value to a collection-key pair.
func (c *Client) PutRaw(collection string, key string, value io.Reader) (*Path, error) {
	return c.doPut(&Path{Collection: collection, Key: key}, nil, value)
}

// Store a value to a collection-key pair if the path's ref value is the latest.
func (c *Client) PutIfUnmodified(path *Path, value interface{}) (*Path, error) {
	reader, writer := io.Pipe()
	encoder := json.NewEncoder(writer)

	go func() { writer.CloseWithError(encoder.Encode(value)) }()
	return c.PutIfUnmodifiedRaw(path, reader)
}

// Store a value to a collection-key pair if the path's ref value is the latest.
func (c *Client) PutIfUnmodifiedRaw(path *Path, value io.Reader) (*Path, error) {
	headers := map[string]string{
		"If-Match": `"` + path.Ref + `"`,
	}

	return c.doPut(path, headers, value)
}

// Store a value to a collection-key pair if it doesn't already hold a value.
func (c *Client) PutIfAbsent(collection, key string, value interface{}) (*Path, error) {
	reader, writer := io.Pipe()
	encoder := json.NewEncoder(writer)

	go func() { writer.CloseWithError(encoder.Encode(value)) }()
	return c.PutIfAbsentRaw(collection, key, reader)
}

// Store a value to a collection-key pair if it doesn't already hold a value.
func (c *Client) PutIfAbsentRaw(collection, key string, value io.Reader) (*Path, error) {
	headers := map[string]string{
		"If-None-Match": "\"*\"",
	}

	return c.doPut(&Path{Collection: collection, Key: key}, headers, value)
}

// Execute a key/value Put.
func (c *Client) doPut(path *Path, headers map[string]string, value io.Reader) (*Path, error) {
	resp, err := c.doRequest("PUT", path.trailingPutURI(), headers, value)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, newError(resp)
	}

	ref := ""
	parts := strings.SplitAfter(resp.Header.Get("Location"), "/")
	if len(parts) >= 6 {
		ref = parts[5]
	} else {
		return nil, fmt.Errorf("Missing ref component: %s", resp.Header.Get("Location"))
	}

	return &Path{
		Collection: path.Collection,
		Key:        path.Key,
		Ref:        ref,
	}, err
}

// Delete the value held at a collection-key pair.
func (c *Client) Delete(collection, key string) error {
	return c.doDelete(collection+"/"+key, nil)
}

// Delete the value held at a collection-key par if the path's ref value is the
// latest.
func (c *Client) DeleteIfUnmodified(path *Path) error {
	headers := map[string]string{
		"If-Match": `"` + path.Ref + `"`,
	}

	return c.doDelete(path.trailingPutURI(), headers)
}

// Delete the current and all previous values from a collection-key pair.
func (c *Client) Purge(collection, key string) error {
	return c.doDelete(collection+"/"+key+"?purge=true", nil)
}

// Delete a collection.
func (c *Client) DeleteCollection(collection string) error {
	return c.doDelete(collection+"?force=true", nil)
}

// Execute delete
func (c *Client) doDelete(trailingUri string, headers map[string]string) error {
	resp, err := c.doRequest("DELETE", trailingUri, headers, nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return newError(resp)
	}

	return nil
}

// List the values in a collection in key order with the specified page size.
func (c *Client) List(collection string, limit int) (*KVResults, error) {
	queryVariables := url.Values{
		"limit": []string{strconv.Itoa(limit)},
	}

	trailingUri := collection + "?" + queryVariables.Encode()

	return c.doList(trailingUri)
}

// List the values in a collection in key order with the specified page size
// that come after the specified key.
func (c *Client) ListAfter(collection, after string, limit int) (*KVResults, error) {
	queryVariables := url.Values{
		"limit":    []string{strconv.Itoa(limit)},
		"afterKey": []string{after},
	}

	trailingUri := collection + "?" + queryVariables.Encode()

	return c.doList(trailingUri)
}

// List the values in a collection in key order with the specified page size
// starting with the specified key.
func (c *Client) ListStart(collection, start string, limit int) (*KVResults, error) {
	queryVariables := url.Values{
		"limit":    []string{strconv.Itoa(limit)},
		"startKey": []string{start},
	}

	trailingUri := collection + "?" + queryVariables.Encode()

	return c.doList(trailingUri)
}

// List the values in a collection within a given range of keys, starting with the
// specified key and stopping at the end key
func (c *Client) ListRange(collection, start, end string, limit int) (*KVResults, error) {
	queryVariables := url.Values{
		"limit":    []string{strconv.Itoa(limit)},
		"startKey": []string{start},
		"endKey":   []string{end},
	}

	trailingUri := collection + "?" + queryVariables.Encode()

	return c.doList(trailingUri)
}

// Get the page of key/value list results that follow that provided set.
func (c *Client) ListGetNext(results *KVResults) (*KVResults, error) {
	return c.doList(results.Next[4:])
}

// Execute a key/value list operation.
func (c *Client) doList(trailingUri string) (*KVResults, error) {
	resp, err := c.doRequest("GET", trailingUri, nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, newError(resp)
	}

	decoder := json.NewDecoder(resp.Body)
	result := new(KVResults)
	if err := decoder.Decode(result); err != nil {
		return result, err
	}

	return result, nil
}

// Check if there is a subsequent page of key/value list results.
func (r *KVResults) HasNext() bool {
	return r.Next != ""
}

// Marshall the value of a KVResult into the provided object.
func (r *KVResult) Value(value interface{}) error {
	return json.Unmarshal(r.RawValue, value)
}

// Returns the trailing URI part for a GET request.
func (p *Path) trailingGetURI() string {
	if p.Ref != "" {
		return p.Collection + "/" + p.Key + "/refs/" + p.Ref
	}
	return p.Collection + "/" + p.Key
}

// Returns the trailing URI part for a PUT request.
func (p *Path) trailingPutURI() string {
	return p.Collection + "/" + p.Key
}
