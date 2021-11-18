package client

// Option is used to configure Service.
type Option func(service *Scheduler)

// WithToken provides scheduler API access token.
func WithToken(token string) Option {
	return func(s *Scheduler) {
		s.accessToken = token
	}
}
