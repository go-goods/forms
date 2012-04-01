package forms

import (
	"errors"
	"net/url"
	"testing"
)

var (
	load_error = errors.New("load_error")

	error_loader LoaderFunc = func(in map[string]interface{}) (out interface{}, errs map[string]error, err error) {
		err = load_error
		return
	}

	int_loader LoaderFunc = func(in map[string]interface{}) (out interface{}, errs map[string]error, err error) {
		out = 2
		return
	}
)

func fatal_loader(t *testing.T) LoaderFunc {
	return func(in map[string]interface{}) (out interface{}, errs map[string]error, err error) {
		t.Fail()
		return
	}
}

func TestLoaderNotCalledOnConverterError(t *testing.T) {
	f := &Form{
		Fields: []Field{
			{Name: "foo", Converter: error_converter},
		},
		Loader: fatal_loader(t),
	}
	f.Load(create_req(url.Values{
		"foo": {"bar"},
	}))
}

func TestLoaderNotCalledOnValidatorError(t *testing.T) {
	f := &Form{
		Fields: []Field{
			{Name: "foo", Validators: []Validator{error_validator}},
		},
		Loader: fatal_loader(t),
	}
	f.Load(create_req(url.Values{
		"foo": {"bar"},
	}))
}

func TestLoaderReturnsValue(t *testing.T) {
	f := &Form{
		Loader: int_loader,
	}
	res := f.Load(create_req(nil))
	if ex, ok := res.Value.(int); !ok || ex != 2 {
		t.Fatal("Expected %v. Got %v", 2, ex)
	}
}
