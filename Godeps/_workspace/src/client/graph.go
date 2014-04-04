// Copyright 2014, Orchestrate.IO, Inc.

package client

import (
	"encoding/json"
	"fmt"
	"strings"
)

type GraphResults struct {
	Count   uint64        `json:"count"`
	Results []GraphResult `json:"results"`
}

type GraphResult struct {
	Path  ResultPath             `json:"path"`
	Key   string                 `json:"key"`
	Ref   string                 `json:"ref"`
	Value map[string]interface{} `json:"value"`
}

func (client *Client) GetRelations(collection string, key string, hops []string) (*GraphResults, error) {
	relationsPath := strings.Join(hops, "/")

	resp, err := client.doRequest("GET", fmt.Sprintf("%v/%v/relations/%v", collection, key, relationsPath), nil)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, newError(resp)
	}

	decoder := json.NewDecoder(resp.Body)
	result := new(GraphResults)
	if err := decoder.Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (client *Client) PutRelation(sourceCollection string, sourceKey string, kind string, sinkCollection string, sinkKey string) error {
	resp, err := client.doRequest("PUT", fmt.Sprintf("%v/%v/relation/%v/%v/%v", sourceCollection, sourceKey, kind, sinkCollection, sinkKey), nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return newError(resp)
	}
	return nil
}
