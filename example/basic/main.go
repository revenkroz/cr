package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/revenkroz/cr/v2"
	"github.com/revenkroz/cr/v2/example/basic/command/mydomain"
	"github.com/revenkroz/cr/v2/middleware"
)

type Logger struct {
	logger *log.Logger
}

func (l *Logger) Logf(format string, args ...interface{}) {
	l.logger.Printf("[CustomLogger] "+format, args...)
}

func main() {
	logger := &Logger{
		logger: log.New(os.Stderr, "", log.LUTC),
	}

	r := cr.New(
		cr.WithLogger(logger),
		cr.WithMiddleware(middleware.Logger(cr.NewStdLogger())),
	)
	// Register handler
	err := r.Register(&mydomain.Echo{})
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	// JSON params, good data
	mappedRequest := &cr.Command{
		Name:   "MyDomain.Echo",
		Params: json.RawMessage(`{"a": 1, "b": 2}`),
	}
	// JSON params, bad data
	mappedRequest2 := &cr.Command{
		ID:     1,
		Name:   "MyDomain.Echo",
		Params: json.RawMessage(`{"test": 1}`),
	}
	// JSON params, b is 0
	mappedRequest3 := &cr.Command{
		ID:     2,
		Name:   "MyDomain.Echo",
		Params: json.RawMessage(`{"a": 1, "b": 0}`),
	}

	commands := []*cr.Command{
		mappedRequest,
		mappedRequest2,
		mappedRequest3,
	}

	resp := r.Run(cr.NewContext(), commands, true)

	for _, res := range resp {
		if res.Error != nil {
			log.Printf("error: %s", res.Error.Error())
		} else {
			log.Printf("result: %+v", res.Result)
		}
	}
}
