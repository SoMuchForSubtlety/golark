package golark

import "log"

// This is the most basic way to make a request.
// It will request the person with ID "pers_123" with all it's fields.
func ExampleRequest() {
	type person struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	var p person

	r := NewRequest("https://test.com/api/", "person", "pers_123")
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

	r := NewRequest("https://test.com/api/", "person", "")
	err := r.Execute(&people)
	if err != nil {
		log.Fatal(err)
	}
}

// AddField lets you limit requests to certain fields, this can speed them up significantly.
func ExampleRequest_AddField() {
	NewRequest("https://test.com/api/", "person", "pers_123").
		AddField(NewField("first_name")).
		AddField(NewField("last_name"))
}

func ExampleRequest_Expand() {
	NewRequest("https://test.com/api/", "person", "pers_123").
		Expand(NewField("team_url"))
}

func ExampleRequest_OrderBy() {
	NewRequest("https://test.com/api/", "person", "").
		OrderBy(NewField("first_name"))
}

// Filters are helpful to search for onjects with knowing their ID
func ExampleRequest_WithFilter() {
	NewRequest("https://test.com/api/", "person", "").
		WithFilter("first_name", NewFilter(Equals, "Bob"))
}

func ExampleRequest_WithFilter_greater_than() {
	NewRequest("https://test.com/api/", "person", "").
		WithFilter("salary", NewFilter(GreaterThan, "10000"))
}

// You can use comma separated lists to query multiple objects at once.
func ExampleRequest_WithFilter_multiple() {
	NewRequest("https://test.com/api/", "person", "").
		WithFilter("first_name", NewFilter(Equals, "Bob,Lucas,Sue"))
}

func ExampleField_WithSubField() {
	NewRequest("https://test.com/api/", "person", "").
		AddField(NewField("first_name")).
		AddField(NewField("team_url").
			WithSubField(NewField("name")).
			WithSubField(NewField("nation_url").
				WithSubField(NewField("name"))))
}
