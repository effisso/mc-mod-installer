package mc

import "net/http"

// HTTPGet is the interface for making HTTP GET requests abstractly
type HTTPGet interface {
	Get(url string) (*http.Response, error)
}

// HTTPClient contains interfaces for interacting with HTTP web services
type HTTPClient struct {
	Getter HTTPGet
}

// NewHTTPClient uses a live http.Client to make connections
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		Getter: &http.Client{
			CheckRedirect: CheckRedirect,
		},
	}
}

// CheckRedirect makes redirects are followed. Only exported for testing
func CheckRedirect(r *http.Request, _ []*http.Request) error {
	r.URL.Opaque = r.URL.Path
	return nil
}
