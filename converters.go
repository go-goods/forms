package forms

import (
	"errors"
	"strconv"
	"time"
)

const (
	time_format = "2006-01-02"
)

var (
	//Some built in converters that return the type they suggest with any errors.
	IntConverter     ConverterFunc = int_converter
	Float64Converter ConverterFunc = float64_converter
	Float32Converter ConverterFunc = float32_converter
	TimeConverter    ConverterFunc = time_converter
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

func int_converter(in string) (out interface{}, err error) {
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

func float64_converter(in string) (out interface{}, err error) {
	out, err = strconv.ParseFloat(in, 64)

	if numerr, ok := err.(*strconv.NumError); ok && err != nil {
		err = make_human_readable(numerr)
		return
	}
	return
}

func float32_converter(in string) (out interface{}, err error) {
	f, err := strconv.ParseFloat(in, 32)

	if numerr, ok := err.(*strconv.NumError); ok && err != nil {
		err = make_human_readable(numerr)
		return
	}
	out = float32(f)
	return
}

func time_converter(in string) (out interface{}, err error) {
	//parse the input
	t, err := time.Parse(time_format, in)

	if err != nil {
		return nil, errors.New("Invalid time: YYYY-MM-DD")
	}

	//set our output
	out = time.Time(t)
	return
}
