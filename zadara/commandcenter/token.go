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

func (t *addTokenHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Token", t.token)

	res, err := t.T.RoundTrip(req)
	if err != nil {
		return res, fmt.Errorf("error making request: %w", err)
	}

	return res, nil
}

func newAddTokenHeaderTransport(
	roundTripper http.RoundTripper,
	token string,
) *addTokenHeaderTransport {
	if roundTripper == nil {
		roundTripper = http.DefaultTransport
	}

	return &addTokenHeaderTransport{T: roundTripper, token: token}
}
