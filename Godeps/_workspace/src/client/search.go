package client

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
)

type SearchResults struct {
	Count      uint64         `json:"count"`
	TotalCount uint64         `json:"total_count"`
	Results    []SearchResult `json:"results"`
}

type SearchResult struct {
	Path  ResultPath             `json:"path"`
	Score float64                `json:"score"`
	Value map[string]interface{} `json:"value"`
}

type ResultPath struct {
	Collection string `json:"collection"`
	Key        string `json:"key"`
	Ref        string `json:"ref"`
}

func (client Client) Search(collection string, query string, limit int, offset int) (*SearchResults, error) {
	queryVariables := url.Values{
		"query":  []string{query},
		"limit":  []string{strconv.Itoa(limit)},
		"offset": []string{strconv.Itoa(offset)},
	}

	resp, err := client.doRequest("GET", collection+"?"+queryVariables.Encode(), nil)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, newError(resp)
	}

	decoder := json.NewDecoder(resp.Body)
	results := new(SearchResults)
	err = decoder.Decode(results)

	if err != nil {
		log.Fatal(err)
	}

	return results, err
}
