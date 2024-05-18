# CR — Go Command Runner

CR is a simple and extensible command runner for Go that allows you to run commands/actions/procedures in a more convenient way.

## Features

- ✨ Easy to use, almost dependency-free
- 📦 Ability to use middlewares for your commands
- 👌 You can use validators for your commands
- 📩 Stamps to pass any data from middleware to command

## Usage

See the [examples](./example) for more information.

### 1. Create runner

We will use `github.com/rs/zerolog/log` as logger.

```go
package runner

import (
	"github.com/revenkroz/cr/runner"
	"github.com/rs/zerolog/log"
)

type Logger struct {
}

func (l *Logger) Logf(format string, args ...interface{}) {
	log.Info().Msgf(format, args...)
}

func CreateLogger() *Logger {
	return &Logger{}
}

type Runner struct {
	*runner.Runner
}

func NewRunner() *Runner {
	return &Runner{
		Runner: runner.New(
			runner.WithLogger(CreateLogger()),
		),
	}
}
```

### 2. Create command

```go
package mydomain

import (
	"errors"
	"github.com/revenkroz/cr/runner"
)

type EchoArgs struct {
	A int `json:"a"`
}

type EchoResponse struct {
	C int `json:"c"`
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
	quo := &EchoResponse{
		C: args.A * 2,
	}

	return quo, nil
}
```

### 3. Register command

```go
// runnerService := NewRunner()
runnerService.Register(&mydomain.Echo{})
```

### 4. Run

```go
result := runnerService.RunOne(runner.NewContext(), &runner.Command{
    Name: "MyDomain.Echo",
    Params: map[string]interface{}{
        "A": 4,
    },
})
```
