package main

import (
	"bytes"
	"client"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"github.com/hoisie/web"
)

var (
	c = client.NewClient(os.Getenv("ORC_KEY"))
)

func main() {
	port := os.Getenv("PORT")
	log.Printf("Listening on port %v ...", port)
	web.Get("/", search)
	web.Run(":"+port)
}

func search(ctx *web.Context) {
	ctx.ContentType("json")

	query := ctx.Params["query"]

	var limit, offset int64
	var err error

	if limit, err = strconv.ParseInt(ctx.Params["limit"], 10, 32); err != nil {
		limit = 10
}
	if offset, err = strconv.ParseInt(ctx.Params["offset"], 10, 32); err != nil {
		offset = 0
	}

	results, err := c.Search("emails", query, int(limit), int(offset))

	if err != nil {
		status, _ := strconv.ParseInt(err.(*client.OrchestrateError).Status[0:3], 10, 32)
		ctx.Abort(int(status), err.Error())
	}

	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(results)

	ctx.Write(buf.Bytes())
}
