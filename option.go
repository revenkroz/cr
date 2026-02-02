package cr

type Option func(r *Runner)

func WithMiddleware(middleware Middleware) Option {
	return func(r *Runner) {
		r.middlewares = append(r.middlewares, middleware)
	}
}

func WithLogger(logger Logger) Option {
	return func(r *Runner) {
		r.logger = logger
	}
}
