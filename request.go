package golark

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type order string

var errHTTP = errors.New("response has non 200 status code")

const (
	// Ascending is used to sort by ascending
	Ascending = order("")
	// Descending is used to sort by descending
	Descending = order("-")
)

// Request represents a Skylark API request
type Request struct {
	Endpoint         string
	Collection       string
	ID               string
	Context          context.Context
	Client           *http.Client
	fields           map[string]*Field
	additionalFields map[string]string
}

// NewRequest returns a simple request with the given
func NewRequest(endpoint, collection, id string) *Request {
	if endpoint[len(endpoint)-1] != '/' {
		endpoint += "/"
	}

	return &Request{
		Collection:       collection,
		Endpoint:         endpoint,
		ID:               id,
		fields:           make(map[string]*Field),
		additionalFields: make(map[string]string),
		Context:          context.Background(),
		Client:           http.DefaultClient,
	}
}

// WithClient lets you set a client that will be used for HTTP requests
// by default http.DefaultCLient is used
func (r *Request) WithClient(client *http.Client) *Request {
	r.Client = client
	return r
}

// AddField adds a field to the request.
// If a request has fields specified it will only return those fields.
func (r *Request) AddField(f *Field) *Request {
	r.fields[f.name] = f
	return r
}

// QueryParams calculates and returns the request's query parameters.
func (r *Request) QueryParams() url.Values {
	v := url.Values{}
	for _, field := range r.fields {
		v = field.apply(v, "")
	}
	for key, value := range r.additionalFields {
		v.Add(key, value)
	}
	return v
}

// OrderBy sorts the response by the given field
func (r *Request) OrderBy(f *Field, order order) *Request {
	r.additionalFields["order"] = string(order) + f.name
	return r
}

// WithFilter allows to filter by a field that is not in the requested response
func (r *Request) WithFilter(fieldName string, filter *Filter) *Request {
	if filter.c != Equals {
		fieldName = fmt.Sprintf("%s__%s", fieldName, filter.c)
	}
	r.additionalFields[fieldName] = filter.value
	return r
}

// Expand expands a field without explicitly listing it as a field to return.
// This is usefult if you want to return all fields.
func (r *Request) Expand(f *Field) *Request {
	f.isExpanded = true
	f.isIncluded = false
	r.AddField(f)
	return r
}

// WithContext set's the context the request will be executed with.
// by default context.Background() is used
func (r *Request) WithContext(ctx context.Context) *Request {
	r.Context = ctx
	return r
}

// Execute executes the request and writes it's results to the value pointed to by v.
func (r *Request) Execute(v interface{}) error {
	url, err := r.ToURL()
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(r.Context, http.MethodGet, url.String(), nil)
	if err != nil {
		return err
	}
	res, err := r.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		message, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errHTTP
		}
		return fmt.Errorf("%s : %w", string(message), errHTTP)
	}

	return json.NewDecoder(res.Body).Decode(v)
}

// ToURL converts the request into a url.URL
func (r *Request) ToURL() (*url.URL, error) {
	temp := r.Endpoint + r.Collection + "/"
	if r.ID != "" {
		temp += r.ID + "/"
	}
	queryParams := r.QueryParams().Encode()
	if queryParams != "" {
		temp += "?" + queryParams
	}
	return url.Parse(temp)
}
