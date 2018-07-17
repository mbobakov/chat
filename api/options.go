package api

type option func(s *server)

// WithLogger setup logger for the serverside. All logging has templae '[LEVEL] Message'
func WithLogger(l logger) option {
	if l == nil {
		return func(s *server) {}
	}
	return func(s *server) {
		s.log = l
	}
}
