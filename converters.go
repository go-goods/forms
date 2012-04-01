package forms

import (
	"errors"
	"strconv"
)

func make_human_readable(numerr *strconv.NumError) (err error) {
	switch numerr.Err {
	case strconv.ErrRange:
		err = errors.New("That number is out of range")
	case strconv.ErrSyntax:
		err = errors.New("That is not a number")
	}
	return
}

var (
	IntConverter ConverterFunc = func(in string) (out interface{}, err error) {
		//parse the input
		i, err := strconv.ParseInt(in, 10, 0)

		//attempt to make the errors more human readable
		if numerr, ok := err.(*strconv.NumError); ok && err != nil {
			err = make_human_readable(numerr)
			return
		}

		//set our output
		out = int(i)
		return
	}

	Float64Converter ConverterFunc = func(in string) (out interface{}, err error) {
		out, err = strconv.ParseFloat(in, 64)

		if numerr, ok := err.(*strconv.NumError); ok && err != nil {
			err = make_human_readable(numerr)
			return
		}
		return
	}

	Float32Converter ConverterFunc = func(in string) (out interface{}, err error) {
		f, err := strconv.ParseFloat(in, 32)

		if numerr, ok := err.(*strconv.NumError); ok && err != nil {
			err = make_human_readable(numerr)
			return
		}
		out = float32(f)
		return
	}
)
