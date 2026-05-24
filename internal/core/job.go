package core

import (
	"net/http"
	"net/url"
	"sync/atomic"
)

type Job struct {
	ID  uint32
	URL *url.URL

	// Body is marshalled payload bytes.
	Method string
	Body   []byte

	Headers http.Header
}

func CreateJobFactory(urlStr string, method string, body []byte, headers http.Header) func() Job {

	url, err := url.Parse(urlStr)
	if err != nil {
		return nil
	}

	var id uint32

	return func() Job {
		return Job{
			ID:  atomic.AddUint32(&id, 1),
			URL: url,

			Method: method,
			Body:   body,

			Headers: headers,
		}
	}
}
