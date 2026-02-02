package mydomain

import (
	"errors"
	"github.com/revenkroz/cr"
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

func (h *Echo) Handler() cr.HandlerFunc {
	return cr.H(h.Echo)
}

func (h *Echo) Echo(ctx cr.Context, args *EchoArgs) (*EchoResponse, error) {
	if args.A == 0 {
		return nil, errors.New("echo zero")
	}

	if args.B == 0 {
		return nil, cr.NewValidationError([]cr.Violation{
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
