orchestrate-go-client
=====================

A golang client for Orchestrate.io

Go Style Documentation:
[http://godoc.org/github.com/orchestrate-io/orchestrate-go-client](http://godoc.org/github.com/orchestrate-io/orchestrate-go-client)]

Usage examples

```go
    c := client.NewClient("Your API Key")

    // Get a value
    value, _ := c.Get("collection", "key")

    // Put a value
    c.Put("collection", "key", strings.NewReader("Some JSON"))

    // Search
    results, _ := c.Search("collection", "A Lucene Query")

    // Get Events
    events, _ := c.GetEvents("collection", "key", "kind")

    // Put Event
    c.PutEvent("collection", "key", "kind", strings.NewReader("Some JSON"))

    // Get Relations
    relations, _ := c.GetRelations("collection", "key", []string{"kind", "kind"})

    // Put Relation
    c.PutRelation("sourceCollection", "sourceKey", "kind", "sinkCollection", "sinkKey")
```
