package middleware

import (
	"github.com/revenkroz/cr"
	"time"
)

func Logger(logger cr.Logger) cr.Middleware {
	return func(handler cr.MiddlewareFunc) cr.MiddlewareFunc {
		return func(ctx cr.Context, cmd *cr.Command) *cr.Result {
			t1 := time.Now().UnixMicro()
			resp := handler(ctx, cmd)
			t2 := time.Now().UnixMicro()

			logger.Logf("call=%s, time=%dμs", cmd.Name, t2-t1)

			return resp
		}
	}
}
