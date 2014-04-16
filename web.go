package main

import (
	"bytes"
	"github.com/orchestrate-io/gorc"
	"encoding/json"
	"github.com/hoisie/web"
	"log"
	"os"
	"strconv"
)

var (
	c = gorc.NewClient(os.Getenv("ORC_KEY"))
)

func main() {
	port := os.Getenv("PORT")
	log.Printf("Listening on port %v ...", port)
	web.Get("/([^/]+/?)", search)
	web.Run(":" + port)
}

func search(ctx *web.Context, collection string) {
	ctx.ContentType("json")
	ctx.SetHeader("Access-Control-Allow-Origin", "*", true)

	query := ctx.Params["query"]

	var limit, offset int64
	var err error

	if limit, err = strconv.ParseInt(ctx.Params["limit"], 10, 32); err != nil {
		limit = 10
	}
	if offset, err = strconv.ParseInt(ctx.Params["offset"], 10, 32); err != nil {
		offset = 0
	}

	results, err := c.Search(collection, query, int(limit), int(offset))

	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)

	if err != nil {
		encoder.Encode(err)
		ctx.WriteHeader(err.(*gorc.OrchestrateError).StatusCode)
	} else {
		encoder.Encode(results)
	}

	ctx.Write(buf.Bytes())
}
