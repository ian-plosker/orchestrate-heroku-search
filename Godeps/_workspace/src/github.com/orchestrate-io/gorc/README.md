gorc
====

A golang client for Orchestrate.io

Supports go 1.1 or later

Go Style Documentation:
[http://godoc.org/github.com/orchestrate-io/gorc](http://godoc.org/github.com/orchestrate-io/gorc)

Usage examples

```go
    // Import the client
    import "github.com/orchestrate-io/gorc"

    // Create a client
    c := gorc.NewClient("Your API Key")

    // Get a value
    result, _ := c.Get("collection", "key")

    // Marshall value into a map
    valueMap := make(map[string]interface{})
    result.Value(&valueMap)

    // Marshall value into a domain type
    domainObject := DomainObject{}
    result.Value(&domainObject)

    // Put a serialized value
    c.PutRaw("collection", "key", strings.NewReader("Some JSON"))

    // Put a interface{} type
    group := Group{Name: "name", Founded: 1990}
    c.Put("collection", "key", group)

    // Search
    results, _ := c.Search("collection", "A Lucene Query", 100, 0)

    // Marshall (search/event/graph) results into an array of domain objects
    var values []DomainObject{} = make([]DomainObject{}, len(results.Results))
    for i, result := range results.Results {
        result.Value(&values[i])
    }

    // Get next page of results
    if results.HasNext() {
        results, err := c.SearchGetNext(results)
    }

    // Get Events
    events, _ := c.GetEvents("collection", "key", "kind")

    // Put Events
    c.PutEvent("collection", "key", "kind", domainObject)
    c.PutEventRaw("collection", "key", "kind", strings.NewReader(serializedJson))

    // Get Relations
    relations, _ := c.GetRelations("collection", "key", []string{"kind", "kind"})

    // Put Relation
    c.PutRelation("sourceCollection", "sourceKey", "kind", "sinkCollection", "sinkKey")
```
