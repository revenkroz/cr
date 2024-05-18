package runner

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Runner struct {
	logger      Logger
	lock        Lock
	handlers    map[string]CommandHandler
	middlewares []Middleware
}

func New(opts ...Option) *Runner {
	s := &Runner{
		logger:   NewVoidLogger(),
		lock:     NewMutexLock(),
		handlers: map[string]CommandHandler{},
	}

	s.Use(opts...)

	return s
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
	r.lock.Lock()
	defer r.lock.Unlock()

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
	r.lock.ReadLock()
	defer r.lock.ReadUnlock()
	_, ok := r.handlers[strings.ToLower(name)]

	return ok
}

func (r *Runner) Get(name string) CommandHandler {
	r.lock.ReadLock()
	defer r.lock.ReadUnlock()
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
) (resp []*Result) {
	wg := sync.WaitGroup{}

	exec := func(req *Command) *Result {
		h := func(ctx Context, req *Command) *Result {
			return r.callMethodMiddleware(ctx, req)
		}
		for _, m := range r.middlewares {
			h = m(h)
		}

		return h(ctx, req)
	}

	if parallel {
		wg.Add(len(commands))
		for _, req := range commands {
			cmd := req

			go func(cmd *Command) {
				defer wg.Done()
				resp = append(resp, exec(cmd))
			}(cmd)
		}
	} else {
		for _, req := range commands {
			resp = append(resp, exec(req))
		}
	}

	if parallel {
		wg.Wait()
	}

	return
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
