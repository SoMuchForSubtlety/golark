package golark_test

import (
	"log"

	"github.com/SoMuchForSubtlety/golark"
)

// This is the most basic way to make a request.
// It will request the person with ID "pers_123" with all it's fields.
func ExampleRequest() {
	type person struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	var p person

	r := golark.NewRequest("https://test.com/api/", "person", "pers_123")
	err := r.Execute(&p)
	if err != nil {
		log.Fatal(err)
	}
}

// If you leave the ID empty all members of the collection will be queried.
// You usually want to avoid this unless you limit the request to certain fields.
func ExampleRequest_all() {
	type person struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	type response struct {
		Objects []person `json:"objects"`
	}

	var people []response

	r := golark.NewRequest("https://test.com/api/", "person", "")
	err := r.Execute(&people)
	if err != nil {
		log.Fatal(err)
	}
}

// AddField lets you limit requests to certain fields, this can speed them up significantly.
func ExampleRequest_AddField() {
	golark.NewRequest("https://test.com/api/", "person", "pers_123").
		AddField(golark.NewField("first_name")).
		AddField(golark.NewField("last_name"))
}

func ExampleRequest_Expand() {
	golark.NewRequest("https://test.com/api/", "person", "pers_123").
		Expand(golark.NewField("team_url"))
}

func ExampleRequest_OrderBy() {
	golark.NewRequest("https://test.com/api/", "person", "").
		OrderBy(golark.NewField("first_name"), golark.Ascending)
}

// Filters are helpful to search for objects with knowing their ID
func ExampleRequest_WithFilter() {
	golark.NewRequest("https://test.com/api/", "person", "").
		WithFilter("first_name", golark.NewFilter(golark.Equals, "Bob"))
}

func ExampleRequest_WithFilter_greater_than() {
	golark.NewRequest("https://test.com/api/", "person", "").
		WithFilter("salary", golark.NewFilter(golark.GreaterThan, "10000"))
}

// You can use comma separated lists to query multiple objects at once.
func ExampleRequest_WithFilter_multiple() {
	golark.NewRequest("https://test.com/api/", "person", "").
		WithFilter("first_name", golark.NewFilter(golark.Equals, "Bob,Lucas,Sue"))
}

func ExampleField_WithSubField() {
	golark.NewRequest("https://test.com/api/", "person", "").
		AddField(golark.NewField("first_name")).
		AddField(golark.NewField("team_url").
			WithSubField(golark.NewField("name")).
			WithSubField(golark.NewField("nation_url").
				WithSubField(golark.NewField("name"))))
}
