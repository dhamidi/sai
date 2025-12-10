package java

import (
	"encoding/json"
	"net/url"
)

type URLString struct {
	url.URL
}

func NewURLString(u *url.URL) URLString {
	if u == nil {
		return URLString{}
	}
	return URLString{URL: *u}
}

func FileURL(path string) URLString {
	return URLString{
		URL: url.URL{
			Scheme: "file",
			Path:   path,
		},
	}
}

func (u URLString) IsZero() bool {
	return u.URL.Scheme == "" && u.URL.Host == "" && u.URL.Path == ""
}

func (u URLString) MarshalJSON() ([]byte, error) {
	if u.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(u.URL.String())
}

func (u *URLString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == nil {
		*u = URLString{}
		return nil
	}
	parsed, err := url.Parse(*s)
	if err != nil {
		return err
	}
	u.URL = *parsed
	return nil
}
