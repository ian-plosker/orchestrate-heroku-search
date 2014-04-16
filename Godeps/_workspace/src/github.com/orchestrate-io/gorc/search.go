// Copyright 2014, Orchestrate.IO, Inc.

package gorc

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// Holds results returned from a Search query.
type SearchResults struct {
	Count      uint64         `json:"count"`
	TotalCount uint64         `json:"total_count"`
	Results    []SearchResult `json:"results"`
	Next       string         `json:"next,omitempty"`
	Prev       string         `json:"prev,omitempty"`
}

// An individual search result.
type SearchResult struct {
	Path     Path            `json:"path"`
	Score    float64         `json:"score"`
	RawValue json.RawMessage `json:"value"`
}

// Search a collection with a Lucene Query Parser Syntax Query
// (http://lucene.apache.org/core/4_5_1/queryparser/org/apache/lucene/queryparser/classic/package-summary.html#Overview)
// and with a specified size limit and offset.
func (c *Client) Search(collection, query string, limit, offset int) (*SearchResults, error) {
	queryVariables := url.Values{
		"query":  []string{query},
		"limit":  []string{strconv.Itoa(limit)},
		"offset": []string{strconv.Itoa(offset)},
	}

	trailingUri := collection + "?" + queryVariables.Encode()

	return c.doSearch(trailingUri)
}

// Get the page of search results that follow that provided set.
func (c *Client) SearchGetNext(results *SearchResults) (*SearchResults, error) {
	return c.doSearch(results.Next[4:])
}

// Get the page of search results that precede that provided set.
func (c *Client) SearchGetPrev(results *SearchResults) (*SearchResults, error) {
	return c.doSearch(results.Prev[4:])
}

// Execute a search request.
func (c *Client) doSearch(trailingUri string) (*SearchResults, error) {
	resp, err := c.doRequest("GET", trailingUri, nil, nil)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, newError(resp)
	}

	decoder := json.NewDecoder(resp.Body)
	result := new(SearchResults)
	if err := decoder.Decode(result); err != nil {
		return result, err
	}

	return result, nil
}

// Check if there is a subsequent page of search results.
func (r *SearchResults) HasNext() bool {
	return r.Next != ""
}

// Check if there is a previous page of search results.
func (r *SearchResults) HasPrev() bool {
	return r.Prev != ""
}

// Marshall the value of a SearchResult into the provided object.
func (r *SearchResult) Value(value interface{}) error {
	return json.Unmarshal(r.RawValue, value)
}
