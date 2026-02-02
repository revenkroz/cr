package cr

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Runner struct {
	mu         sync.RWMutex
	logger     Logger
	handlers   map[string]CommandHandler
	middlewares []Middleware
}

func New(opts ...Option) *Runner {
	r := &Runner{
		logger:   NewVoidLogger(),
		handlers: map[string]CommandHandler{},
	}

	r.Use(opts...)

	return r
}

func (r *Runner) Use(opts ...Option) {
	for _, opt := range opts {
		opt(r)
	}
}

func (r *Runner) MustRegister(handler CommandHandler) {
	if err := r.Register(handler); err != nil {
		panic(err)
	}
}

func (r *Runner) Register(handler CommandHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Logf("Register method %s", handler.Name())
	if _, ok := r.handlers[strings.ToLower(handler.Name())]; ok {
		return fmt.Errorf("method %s already registered", handler.Name())
	}

	r.handlers[strings.ToLower(handler.Name())] = handler

	return nil
}

func (r *Runner) RegisterMany(handlers []CommandHandler) error {
	for _, h := range handlers {
		err := r.Register(h)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Runner) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.handlers[strings.ToLower(name)]

	return ok
}

func (r *Runner) Get(name string) CommandHandler {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.handlers[strings.ToLower(name)]
	if !ok {
		return nil
	}

	return h
}

func (r *Runner) RunOne(
	ctx Context,
	req *Command,
) *Result {
	responses := r.Run(ctx, []*Command{req}, false)
	if len(responses) == 0 {
		return nil
	}

	return responses[0]
}

func (r *Runner) Run(
	ctx Context,
	commands []*Command,
	parallel bool,
) []*Result {
	exec := func(req *Command) *Result {
		h := func(ctx Context, req *Command) *Result {
			return r.callMethodMiddleware(ctx, req)
		}
		for _, m := range r.middlewares {
			h = m(h)
		}

		return h(ctx, req)
	}

	resp := make([]*Result, len(commands))

	if parallel {
		wg := sync.WaitGroup{}
		wg.Add(len(commands))
		for i, req := range commands {
			i, cmd := i, req
			go func() {
				defer wg.Done()
				resp[i] = exec(cmd)
			}()
		}
		wg.Wait()
	} else {
		for i, req := range commands {
			resp[i] = exec(req)
		}
	}

	return resp
}

func (r *Runner) callMethodMiddleware(ctx Context, cmd *Command) *Result {
	h := r.Get(cmd.Name)
	if h == nil {
		return ErrorResponse(cmd.ID, NewNotFoundError())
	}

	resp, err := h.Handler()(ctx, cmd.Params)
	if err != nil {
		r.logger.Logf("execution error: %v", err)

		if errors.As(err, &Error{}) {
			return ErrorResponse(cmd.ID, err.(Error))
		}

		return ErrorResponse(cmd.ID, NewOtherError(err.Error()))
	}

	return ResultResponse(cmd.ID, resp)
}
