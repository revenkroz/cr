package runner

type Stamp interface {
	GetName() string
}

func GetStamp[T Stamp](ctx Context, name string) *T {
	stamp, ok := ctx.Stamps[name]
	if !ok {
		return nil
	}

	switch stamp := stamp.(type) {
	case T:
		return &stamp
	default:
		return nil
	}
}
