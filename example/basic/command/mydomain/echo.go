package mydomain

import (
	"errors"
	"github.com/revenkroz/cr/runner"
)

type EchoArgs struct {
	A int `json:"a"`
	B int `json:"b"`
}

type EchoResponse struct {
	C int `json:"c"`
	D int `json:"d"`
}

type Echo struct {
}

func (h *Echo) Name() string {
	return "MyDomain.Echo"
}

func (h *Echo) Handler() runner.HandlerFunc {
	return runner.H(h.Echo)
}

func (h *Echo) Echo(ctx runner.Context, args *EchoArgs) (*EchoResponse, error) {
	if args.A == 0 {
		return nil, errors.New("echo zero")
	}

	if args.B == 0 {
		return nil, runner.NewValidationError([]runner.Violation{
			{
				Field: "B",
				Error: "B cannot be zero",
			},
		})
	}

	quo := &EchoResponse{
		C: args.A,
		D: args.B,
	}

	return quo, nil
}
