package main

import (
	"encoding/json"
	"github.com/revenkroz/cr/example/basic/command/mydomain"
	"github.com/revenkroz/cr/middleware"
	"github.com/revenkroz/cr/runner"
	"log"
	"os"
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

	runnerService := runner.New(
		runner.WithLogger(logger),
		runner.WithMiddleware(middleware.Logger(runner.NewStdLogger())),
	)
	// Register handler
	err := runnerService.Register(&mydomain.Echo{})
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	// JSON params, good data
	mappedRequest := &runner.Command{
		Name:   "MyDomain.Echo",
		Params: json.RawMessage(`{"a": 1, "b": 2}`),
	}
	// JSON params, bad data
	mappedRequest2 := &runner.Command{
		ID:     1,
		Name:   "MyDomain.Echo",
		Params: json.RawMessage(`{"test": 1}`),
	}
	// JSON params, b is 0
	mappedRequest3 := &runner.Command{
		ID:     2,
		Name:   "MyDomain.Echo",
		Params: json.RawMessage(`{"a": 1, "b": 0}`),
	}
	// mapped params
	mappedRequest4 := &runner.Command{
		ID:   3,
		Name: "MyDomain.Echo",
		Params: map[string]interface{}{
			"A": 4,
			"B": 8,
		},
	}

	commands := []*runner.Command{
		mappedRequest,
		mappedRequest2,
		mappedRequest3,
		mappedRequest4,
	}

	resp := runnerService.Run(runner.NewContext(), commands, true)

	for _, r := range resp {
		if r.Error != nil {
			log.Printf("error: %s", r.Error.Error())
		} else {
			log.Printf("result: %+v", r.Result)
		}
	}
}
