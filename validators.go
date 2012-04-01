package forms

import "errors"

var (
	NonemptyValidator ValidatorFunc = func(in string) (out string, err error) {
		out = in
		if len(in) == 0 {
			err = errors.New("Value must be present")
		}
		return
	}
)
