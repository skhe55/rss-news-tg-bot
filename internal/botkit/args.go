package botkit

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

func ParseJSON[T any](src string) (T, error) {
	var args T

	if err := json.Unmarshal([]byte(src), &args); err != nil {
		return *(new(T)), err
	}

	validate := validator.New()
	if err := validate.Struct(args); err != nil {
		return *(new(T)), err
	}

	return args, nil
}
