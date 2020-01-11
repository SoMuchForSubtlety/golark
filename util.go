package client

import "net/url"

func addValue(v url.Values, key string, value string) url.Values {
	values := v.Get(key)
	if values != "" {
		values += ","
	}
	values += value
	v.Set(key, values)
	return v
}
