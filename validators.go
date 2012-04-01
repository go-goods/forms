package forms

import "errors"

var (
	NonemptyValidator ValidatorFunc = nonempty_validator
)

func nonempty_validator(in string) (out string, err error) {
	out = in
	if len(in) == 0 {
		err = errors.New("Value must be present")
	}
	return
}
