package urlbuilder

import (
	"net/http"
	"net/url"
	"strings"
)

// Build URL with query parameters
func BuildWithQueryParams(baseUrl string, params map[string]string) (string, error) {
	u, err := url.Parse(baseUrl)

	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func BaseURL(r *http.Request) string {
	var b strings.Builder

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	b.WriteString(scheme)
	b.WriteString("://")
	b.WriteString(r.Host)
	b.WriteString(r.URL.Path)

	return b.String()
}
