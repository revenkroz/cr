package cr

import "encoding/json"

type HandlerFunc func(Context, interface{}) (interface{}, error)

type CommandHandler interface {
	Name() string
	Handler() HandlerFunc
}

type Command struct {
	ID     int              `json:"id"`
	Name   string           `json:"name"`
	Params json.RawMessage  `json:"params"`
}

type Result struct {
	ID     int         `json:"id,omitempty"`
	Result interface{} `json:"result,omitempty"`
	Error  *Error      `json:"error,omitempty"`
}

func ResultResponse(id int, resp interface{}) *Result {
	return &Result{
		Result: resp,
		ID:     id,
	}
}

func ErrorResponse(id int, err Error) *Result {
	return &Result{
		Error: &err,
		ID:    id,
	}
}

type MiddlewareFunc func(ctx Context, req *Command) *Result

type Middleware func(handler MiddlewareFunc) MiddlewareFunc
