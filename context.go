package cr

type Validator func(params interface{}) []Violation

type Context struct {
	Stamps    map[string]Stamp
	Validator Validator
}

func NewContext(stamps ...Stamp) Context {
	ctx := Context{
		Stamps: make(map[string]Stamp),
	}

	for _, s := range stamps {
		ctx.AddStamp(s)
	}

	return ctx
}

func (c *Context) AddStamp(stamp Stamp) {
	c.Stamps[stamp.GetName()] = stamp
}

func (c *Context) GetStamp(name string) *Stamp {
	if stamp, ok := c.Stamps[name]; ok {
		return &stamp
	}

	return nil
}

func (c *Context) HasStamp(name string) bool {
	return c.GetStamp(name) != nil
}

func (c *Context) SetValidator(validator Validator) {
	c.Validator = validator
}

func (c *Context) Validate(params interface{}) []Violation {
	if c.Validator == nil {
		return nil
	}

	return c.Validator(params)
}

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
