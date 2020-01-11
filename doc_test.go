package client

import "log"

// This is the most basic way to make a request.
// It will request the person with ID "pers_123" with all it's fields.
func RequestExamples() {
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
func RequestExamples_AllPeople() {
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
