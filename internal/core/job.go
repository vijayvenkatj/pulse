package core

import (
	"net/http"
	"net/url"
	"sync/atomic"
)

type Job struct {
	ID  uint32
	URL *url.URL

	// Body is pre-marshalled
	Method string
	Body   []byte

	Headers http.Header
}

func CreateJobFactory(urlStr string, method string, body []byte, headers http.Header) func() Job {

	url, err := url.Parse(urlStr)
	if err != nil {
		return nil
	}

	return func() Job {
		var id uint32
		return Job{
			ID:  atomic.AddUint32(&id, 1),
			URL: url,

			Method: method,
			Body:   body,

			Headers: headers,
		}
	}
}
