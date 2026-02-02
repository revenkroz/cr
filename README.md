# CR — Go Command Runner

CR is a simple and extensible command runner for Go that allows you to run commands/actions/procedures in a more convenient way.

## Features

- Zero dependencies
- Middleware support
- Context validators
- Stamps to pass data from middleware to command
- Parallel command execution
- Type-safe generic helpers

## Installation

```bash
go get -u github.com/revenkroz/cr/v2
```

## Usage

See the [examples](./example) for more information.

### 1. Create runner

```go
import (
    "github.com/revenkroz/cr/v2"
    "github.com/revenkroz/cr/v2/middleware"
)

r := cr.New(
    cr.WithLogger(cr.NewStdLogger()),
    cr.WithMiddleware(middleware.Logger(cr.NewStdLogger())),
)
```

### 2. Create command

```go
package mydomain

import (
    "github.com/revenkroz/cr/v2"
)

type EchoArgs struct {
    A int `json:"a"`
}

type EchoResponse struct {
    C int `json:"c"`
}

type Echo struct{}

func (h *Echo) Name() string {
    return "MyDomain.Echo"
}

func (h *Echo) Handler() cr.HandlerFunc {
    return cr.H(h.Echo)
}

func (h *Echo) Echo(ctx cr.Context, args *EchoArgs) (*EchoResponse, error) {
    return &EchoResponse{C: args.A * 2}, nil
}
```

### 3. Register command

```go
r.MustRegister(&mydomain.Echo{})
```

### 4. Run

```go
result := r.RunOne(cr.NewContext(), &cr.Command{
    Name:   "MyDomain.Echo",
    Params: json.RawMessage(`{"a": 4}`),
})
```

### Parallel execution

```go
commands := []*cr.Command{
    {Name: "MyDomain.Echo", Params: json.RawMessage(`{"a": 1}`)},
    {Name: "MyDomain.Echo", Params: json.RawMessage(`{"a": 2}`)},
}

results := r.Run(cr.NewContext(), commands, true)
```
