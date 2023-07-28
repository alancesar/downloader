package transport

import (
	"net/http"
)

type (
	TokenProvider func() string

	UserAgentRoundTripper struct {
		userAgent string
		next      http.RoundTripper
	}

	AuthorizationRoundTripper struct {
		provider TokenProvider
		next     http.RoundTripper
	}
)

func NewUserAgentRoundTripper(userAgent string, next http.RoundTripper) http.RoundTripper {
	return &UserAgentRoundTripper{
		userAgent: userAgent,
		next:      next,
	}
}

func NewAuthorizationRoundTripper(provider TokenProvider, next http.RoundTripper) http.RoundTripper {
	return &AuthorizationRoundTripper{
		provider: provider,
		next:     next,
	}
}

func (a UserAgentRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	defer closeBody(r)

	newRequest := cloneRequest(r)
	newRequest.Header.Add("User-Agent", a.userAgent)
	return a.next.RoundTrip(newRequest)
}

func (p *AuthorizationRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	defer closeBody(r)

	authorization := p.provider()
	newRequest := cloneRequest(r)
	newRequest.Header.Add("Authorization", authorization)
	return p.next.RoundTrip(newRequest)
}

func cloneRequest(request *http.Request) *http.Request {
	newRequest := new(http.Request)
	*newRequest = *request

	newRequest.Header = make(http.Header, len(request.Header))
	for k, v := range request.Header {
		newRequest.Header[k] = append([]string(nil), v...)
	}

	return newRequest
}

func closeBody(r *http.Request) {
	if r.Body != nil {
		_ = r.Body.Close()
	}
}
