package core

import (
	"net/http"
	"net/url"
)

type Job struct {
	ID  uint32
	URL *url.URL

	// Body is pre-marshalled
	Method string
	Body   []byte

	Headers http.Header
}
