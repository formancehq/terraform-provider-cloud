package pkg

import (
	"net/http"
)

type Creds interface {
	ClientId() string
	ClientSecret() string
	Endpoint() string
	UserAgent() string
}

type RoundTripperFn func(r *http.Request) (*http.Response, error)

func (fn RoundTripperFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
func NewTransport(rt http.RoundTripper, tp TokenProviderImpl) RoundTripperFn {
	return func(r *http.Request) (*http.Response, error) {
		token, err := tp.RefreshToken(r.Context())
		if err != nil {
			return nil, err
		}
		r.Header.Set("Authorization", "Bearer "+token.AccessToken)
		return rt.RoundTrip(r)
	}
}
