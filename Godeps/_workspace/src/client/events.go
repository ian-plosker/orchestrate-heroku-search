package client

import (
	"encoding/json"
	"io"
	"log"
)

type EventResults struct {
	Count   uint64  `json:"count"`
	Results []Event `json:"results"`
}

type Event struct {
	Timestamp uint64                 `json:"timestamp"`
	Value     map[string]interface{} `json:"value"`
}

func (client Client) GetEvents(collection string, key string, kind string) (*EventResults, error) {
	resp, err := client.doRequest("GET", collection+"/"+key+"/events/"+kind, nil)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, newError(resp)
	}

	decoder := json.NewDecoder(resp.Body)
	results := new(EventResults)
	err = decoder.Decode(results)

	return results, err
}

func (client Client) PutEvent(collection string, key string, kind string, value io.Reader) error {
	resp, err := client.doRequest("PUT", collection+"/"+key+"/events/"+kind, value)

	if err != nil {
		log.Fatal(err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		err = newError(resp)
	}

	return err
}
