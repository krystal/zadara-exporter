package commandcenter

import (
	"fmt"
	"net/http"
)

type (
	// addTokenHeaderTransport represents a transport that adds a token header to the request.
	addTokenHeaderTransport struct {
		T     http.RoundTripper
		token string
	}
)

// RoundTrip executes a single HTTP transaction, adding the X-Token header to the request.
// It returns the response received from the server or an error if the request fails.
func (t *addTokenHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Token", t.token)

	res, err := t.T.RoundTrip(req)
	if err != nil {
		return res, fmt.Errorf("error making request: %w", err)
	}

	return res, nil
}

// newAddTokenHeaderTransport creates a new transport that adds a token header to each request.
// If the provided roundTripper is nil, it defaults to http.DefaultTransport.
func newAddTokenHeaderTransport(
	roundTripper http.RoundTripper,
	token string,
) *addTokenHeaderTransport {
	if roundTripper == nil {
		roundTripper = http.DefaultTransport
	}

	return &addTokenHeaderTransport{T: roundTripper, token: token}
}
