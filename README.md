[![GoDoc](https://godoc.org/github.com/SoMuchForSubtlety/golark?status.svg)](https://godoc.org/github.com/SoMuchForSubtlety/golark)
[![Go Report Card](https://goreportcard.com/badge/github.com/SoMuchForSubtlety/golark)](https://goreportcard.com/report/github.com/SoMuchForSubtlety/golark)
![](https://github.com/SoMuchForSubtlety/golark/workflows/Test/badge.svg)

# golark

Golark makes it easy to build Skylark API requests in golang.

```go
package main

import (
	"github.com/SoMuchForSubtlety/golark"
)

type episode struct {
	Title        string   `json:"title"`
	Subtitle     string   `json:"subtitle"`
	UID          string   `json:"uid"`
	DataSourceID string   `json:"data_source_id"`
	Items        []string `json:"items"`
}

func main() {
	var ep episode

	// request an object
	golark.NewRequest("https://test.com/api/", "episodes", "ep_123").
		Execute(&ep)

	// request an object with only certain fields
	golark.NewRequest("https://test.com/api/", "episodes", "ep_123").
		AddField(golark.NewField("title")).
		AddField(golark.NewField("subtitle")).
		AddField(golark.NewField("uid")).
		Execute(&ep)

	type container struct {
		Objects []episode `json:"objects"`
	}

	var eps container

	// request all members of a collection
	golark.NewRequest("https://test.com/api/", "episodes", "").
		Execute(&eps)

	// request all members of a collection with certain properties
	golark.NewRequest("https://test.com/api/", "episodes", "").
		WithFilter("title", golark.NewFilter(golark.Equals, "test episode title")).
		Execute(&eps)
}
```