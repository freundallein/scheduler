package apiserv

// Option is used to configure Service.
type Option func(service *Service)

// WithPort configures a port used for Service.
func WithPort(port string) Option {
	return func(s *Service) {
		s.Port = port
	}
}
