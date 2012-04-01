package forms

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

var (
	validator_error = errors.New("validator_error")

	error_validator ValidatorFunc = func(in string) (out string, err error) {
		err = validator_error
		return
	}
	pass_validator ValidatorFunc = func(in string) (out string, err error) {
		out = in
		return
	}
)

func fatal_validator(t *testing.T) ValidatorFunc {
	return func(in string) (out string, err error) {
		t.Fail()
		out = in
		return
	}
}

func change_to(out string) ValidatorFunc {
	return func(in string) (o string, err error) {
		o = out
		return
	}
}

func create_req(values url.Values) (req *http.Request) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		panic(err)
	}
	req.Form = values
	return
}

func TestValidatorFails(t *testing.T) {
	f := &Form{
		Fields: []Field{
			{
				Name:       "foo",
				Validators: []Validator{error_validator},
			},
		},
	}
	res := f.Load(create_req(url.Values{
		"foo": {"bar"},
	}))

	if ex, ok := res.Errors["foo"]; !ok || ex != validator_error {
		t.Fatalf("Expected %v. Got %v", validator_error, ex)
	}
	if ex, ok := res.Values["foo"]; !ok || ex != "bar" {
		t.Fatalf("Expected %v. Got %v", "bar", ex)
	}
}

func TestValidatorFailsChain(t *testing.T) {
	f := &Form{
		Fields: []Field{
			{
				Name:       "foo",
				Validators: []Validator{pass_validator, pass_validator, error_validator},
			},
		},
	}
	res := f.Load(create_req(url.Values{
		"foo": {"bar"},
	}))

	if ex, ok := res.Errors["foo"]; !ok || ex != validator_error {
		t.Fatalf("Expected %v. Got %v", validator_error, ex)
	}
	if ex, ok := res.Values["foo"]; !ok || ex != "bar" {
		t.Fatalf("Expected %v. Got %v", "bar", ex)
	}
}

func TestValidatorChanges(t *testing.T) {
	f := &Form{
		Fields: []Field{
			{
				Name:       "foo",
				Validators: []Validator{change_to("baz")},
			},
		},
	}
	res := f.Load(create_req(url.Values{
		"foo": {"bar"},
	}))
	if ex, ok := res.Errors["foo"]; ok || ex != nil {
		t.Fatalf("Expected %v. Got %v", nil, ex)
	}
	rval := res.Value.(map[string]interface{})
	if ex, ok := rval["foo"]; !ok || ex != "baz" {
		t.Fatalf("Expected %v. Got %v", "baz", ex)
	}
}

func TestValidatorEarlyExit(t *testing.T) {
	f := &Form{
		Fields: []Field{
			{
				Name:       "foo",
				Validators: []Validator{error_validator, fatal_validator(t)},
			},
		},
	}
	f.Load(create_req(url.Values{
		"foo": {"bar"},
	}))
}

func TestValidatorEmpty(t *testing.T) {
	f := &Form{
		Fields: []Field{
			{Name: "foo"},
		},
	}
	res := f.Load(create_req(url.Values{
		"foo": {"bar"},
	}))
	if ex, ok := res.Errors["foo"]; ok || ex != nil {
		t.Fatalf("Expected %v. Got %v", nil, ex)
	}
	rval := res.Value.(map[string]interface{})
	if ex, ok := rval["foo"]; !ok || ex != "bar" {
		t.Fatalf("Expected %v. Got %v", "bar", ex)
	}
}

func TestValidatorsNonemptyValidator(t *testing.T) {
	f := &Form{
		Fields: []Field{
			{Name: "foo", Validators: []Validator{NonemptyValidator}},
		},
	}
	res := f.Load(create_req(url.Values{
		"foo": {"bar"},
	}))
	if ex, ok := res.Errors["foo"]; ok || ex != nil {
		t.Fatalf("Expected %v. Got %v", nil, ex)
	}

	res = f.Load(create_req(url.Values{
		"foo": {""},
	}))
	if ex, ok := res.Errors["foo"]; !ok || ex == nil {
		t.Fatalf("Expected %v. Got %v", "not nil", ex)
	}

	res = f.Load(create_req(nil))
	if ex, ok := res.Errors["foo"]; !ok || ex == nil {
		t.Fatalf("Expected %v. Got %v", "not nil", ex)
	}
}
