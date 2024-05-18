package middleware

import (
	"github.com/revenkroz/cr/runner"
	"time"
)

func Logger(logger runner.Logger) runner.Middleware {
	return func(handler runner.MiddlewareFunc) runner.MiddlewareFunc {
		return func(ctx runner.Context, cmd *runner.Command) *runner.Result {
			t1 := time.Now().UnixMicro()
			resp := handler(ctx, cmd)
			t2 := time.Now().UnixMicro()

			logger.Logf("call=%s, time=%dμs", cmd.Name, t2-t1)

			return resp
		}
	}
}
