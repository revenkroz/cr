package runner

type Option func(s *Runner)

func WithMiddleware(mw Middleware) Option {
	return func(s *Runner) {
		s.middlewares = append(s.middlewares, mw)
	}
}

func WithLogger(l Logger) Option {
	return func(s *Runner) {
		s.logger = l
	}
}

func WithLock(l Lock) Option {
	return func(s *Runner) {
		s.lock = l
	}
}
