package runner

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
)

// H is a helper function to create a HandlerFunc from a handler function
func H[CommandParams any, CommandResponse any](handler func(Context, *CommandParams) (CommandResponse, error)) HandlerFunc {
	return func(ctx Context, in interface{}) (interface{}, error) {
		var params *CommandParams
		switch v := in.(type) {
		case *CommandParams:
			params = v
		case map[string]interface{}:
			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				WeaklyTypedInput: true,
				Result:           &params,
			})
			if err != nil {
				return nil, NewParseError()
			}

			if err := decoder.Decode(v); err != nil {
				return nil, NewParseError()
			}
		case json.RawMessage:
			if err := json.Unmarshal(v, &params); err != nil {
				return nil, NewParseError()
			}
		default:
			return nil, NewParseError()
		}

		violations := ctx.Validate(params)
		if violations != nil {
			return nil, NewValidationError(violations)
		}

		resp, err := handler(ctx, params)
		if err != nil {
			return nil, NewOtherError(err.Error())
		}

		return resp, nil
	}
}

// N is a helper function to create a HandlerFunc from a handler function that does not require any parameters
func N[CommandResponse any](handler func(Context) (CommandResponse, error)) HandlerFunc {
	return func(ctx Context, in interface{}) (interface{}, error) {
		resp, err := handler(ctx)
		if err != nil {
			return nil, NewOtherError(err.Error())
		}

		return resp, nil
	}
}
