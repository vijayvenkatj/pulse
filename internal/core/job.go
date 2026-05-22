package core

import (
	"net/http"
	"net/url"
)

type ReqMethod string

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	PATCH  = "PATCH"
)

type Job struct {
	ID  uint32
	URL *url.URL

	// Body is pre-marshalled
	Method ReqMethod
	Body   []byte

	Headers http.Header
}
