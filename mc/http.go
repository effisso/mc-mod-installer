package mc

import (
	"net/http"
	"strings"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
)

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
	client := &http.Client{
		CheckRedirect: CheckRedirect,
	}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)
	return &HTTPClient{
		Getter: client,
	}
}

// CheckRedirect makes redirects are followed. Only exported for testing
func CheckRedirect(r *http.Request, _ []*http.Request) error {
	r.URL.Opaque = strings.ReplaceAll(r.URL.Path, "+", "%2B")
	return nil
}
