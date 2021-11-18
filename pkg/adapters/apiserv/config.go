package apiserv

// Option is used to configure Service.
type Option func(service *Service)

// WithPort configures a port used for Service.
func WithPort(port string) Option {
	return func(s *Service) {
		s.Port = port
	}
}

// WithToken provides scheduler API access token.
func WithToken(token string) Option {
	return func(s *Service) {
		s.Token = token
	}
}

// WithWorkerToken provides worker API access token.
func WithWorkerToken(token string) Option {
	return func(s *Service) {
		s.WorkerToken = token
	}
}
