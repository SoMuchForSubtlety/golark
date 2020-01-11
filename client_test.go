package golark

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

const (
	driverID = "driv_123"
	teamID   = "team_123"
)

func TestExpandOnly(t *testing.T) {
	request := NewRequest("https://test.com/api/", "team", teamID).
		Expand(NewField("nation_url").
			Expand(NewField("eventoccurrence_urls"))).
		Expand(NewField("driver_urls"))

	testURL(request, "https://test.com/api/team/team_123/?fields_to_expand=driver_urls,nation_url,nation_url__eventoccurrence_urls", t)
}

func TestMultipleExpansion(t *testing.T) {
	request := NewRequest("https://test.com/api/", "session-occurrence", "").
		AddField(NewField("channel_urls").
			WithSubField(NewField("self")).
			WithSubField(NewField("name")).
			WithSubField(NewField("driver_urls").
				WithSubField(NewField("driver_racingnumber")).
				WithSubField(NewField("team_url").
					WithSubField(NewField("name")).
					WithSubField(NewField("colour"))))).
		WithFilter("slug", NewFilter(Equals, "test"))

	testURL(request, "https://test.com/api/session-occurrence/?fields=channel_urls,channel_urls__self,channel_urls__name,channel_urls__driver_urls,channel_urls__driver_urls__driver_racingnumber,channel_urls__driver_urls__team_url,channel_urls__driver_urls__team_url__name,channel_urls__driver_urls__team_url__colour&fields_to_expand=channel_urls,channel_urls__driver_urls,channel_urls__driver_urls__team_url&slug=test", t)
}

func TestRequestFilter(t *testing.T) {
	request := NewRequest("https://test.com/api/", "sets", "").
		AddField(NewField("title")).
		AddField(NewField("self")).
		WithFilter("set_type_slug", NewFilter(Equals, "video"))

	testURL(request, "https://test.com/api/sets/?fields=title,self&set_type_slug=video", t)
}

func TestOrder(t *testing.T) {
	year := NewField("year")
	request := NewRequest("https://test.com/api/", "race-season", "").
		AddField(year).
		AddField(NewField("name")).
		AddField(NewField("self")).
		OrderBy(year)

	testURL(request, "https://test.com/api/race-season/?fields=year,name,self&order=year", t)
}

func TestFieldFilter(t *testing.T) {
	request := NewRequest("https://test.com/api/", "race-season", "").
		AddField(NewField("year").
			WithFilter(NewFilter(GreaterThan, "2017"))).
		AddField(NewField("name")).
		AddField(NewField("self"))

	testURL(request, "https://test.com/api/race-season/?fields=year,name,self&year__gt=2017", t)
}

func TestExpandedField(t *testing.T) {
	request := NewRequest("https://test.com/api/", "driver", driverID).
		AddField(NewField("first_name")).
		AddField(NewField("last_name")).
		AddField(NewField("team_url").
			WithSubField(NewField("name")).
			WithSubField(NewField("colour"))).
		AddField(NewField("driver_tla"))

	testURL(request, "https://test.com/api/driver/driv_123/?fields=first_name%2Clast_name%2Cteam_url%2Cteam_url__colour%2Cteam_url__name%2Cdriver_tla&fields_to_expand=team_url", t)
}

func TestAllFields(t *testing.T) {
	request := NewRequest("https://test.com/api/", "driver", driverID)

	testURL(request, "https://test.com/api/driver/driv_123/", t)
}

func testURL(r *Request, expectedURL string, t *testing.T) {
	expected, err := url.Parse(expectedURL)
	if err != nil {
		t.Error("Invalid expected URL:", err)
	}

	actual, err := r.ToURL()
	if err != nil {
		t.Error("Error generating URL:", err)
	}

	if expected.Path != actual.Path {
		t.Error(fmt.Sprintf("incorrect URL path\nexpected: %s\ngot:      %s", expected.Path, actual.Path))
	}

	if expected.Host != actual.Host {
		t.Error(fmt.Sprintf("incorrect host\nexpected: %s\ngot:      %s", expected.Host, actual.Host))
	}

	compareValues(expected.Query(), actual.Query(), t)
}

func compareCSV(expected, actual string, t *testing.T) {
	expectedMap := make(map[string]int)
	actualMap := make(map[string]int)
	for _, expectedElem := range strings.Split(expected, ",") {
		expectedMap[expectedElem]++
	}
	for _, actualElem := range strings.Split(actual, ",") {
		actualMap[actualElem]++
	}

	for expectedKey, expectedVal := range expectedMap {
		if actualMap[expectedKey] != expectedVal {
			t.Error("URL does not contain", expectedKey)
		}
	}

	for actualKey, actualVal := range actualMap {
		if expectedMap[actualKey] != actualVal {
			t.Error("URL contains unexpected value", actualKey)
		}
	}
}

func compareValues(expected, actual url.Values, t *testing.T) {
	for expectedKey, expectedVal := range expected {
		actualVal, ok := actual[expectedKey]
		if !ok {
			t.Error("URL does not contain query param", expectedKey)
		} else {
			compareCSV(expectedVal[0], actualVal[0], t)
		}
	}

	for actualKey := range actual {
		_, ok := expected[actualKey]
		if !ok {
			t.Error("URL contains unexpected query param", actualKey)
		}
	}
}
