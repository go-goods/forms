package forms

import "errors"

var (
	//Some built in validators that do what you would expect.
	NonemptyValidator ValidatorFunc = nonempty_validator
)

func nonempty_validator(in string) (out string, err error) {
	out = in
	if len(in) == 0 {
		err = errors.New("Value must be present")
	}
	return
}
