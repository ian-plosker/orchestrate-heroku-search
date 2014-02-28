package main

import (
	"bytes"
	"client"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var (
	c = client.NewClient("459157d0-b88f-4e7f-8b54-e3fb952e52ec")
)

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", search)
	fmt.Printf("Listening on port %v ...\n", port)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		panic(err)
	}
}

func search(res http.ResponseWriter, req *http.Request) {
	query := req.FormValue("query")

	var limit, offset int64
	var err error

	if limit, err = strconv.ParseInt(req.FormValue("limit"), 10, 32); err != nil {
		limit = 10
}
	if offset, err = strconv.ParseInt(req.FormValue("offset"), 10, 32); err != nil {
		offset = 0
	}

	results, err := c.Search("emails", query, int(limit), int(offset))

	if err != nil {
		status, _ := strconv.ParseInt(err.(*client.OrchestrateError).Status[0:3], 10, 32)

		http.Error(res, err.Error(), int(status))
		return
	}

	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(results)

	res.Header().Set("Content-Type", "application/json")
	res.Write(buf.Bytes())
}
