package forms

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

var (
	always_fail_err = errors.New("always_fail")

	always_fail ValidatorFunc = func(in string) (out string, err error) {
		err = always_fail_err
		return
	}
	always_valid ValidatorFunc = func(in string) (out string, err error) {
		out = in
		return
	}
)

func fail_if_called(t *testing.T) ValidatorFunc {
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
				Validators: []Validator{always_fail},
			},
		},
	}
	res := f.Load(create_req(url.Values{
		"foo": {"bar"},
	}))

	if ex, ok := res.Errors["foo"]; !ok || ex != always_fail_err {
		t.Fatalf("Expected %v. Got %v", always_fail_err, ex)
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
				Validators: []Validator{always_valid, always_valid, always_fail},
			},
		},
	}
	res := f.Load(create_req(url.Values{
		"foo": {"bar"},
	}))

	if ex, ok := res.Errors["foo"]; !ok || ex != always_fail_err {
		t.Fatalf("Expected %v. Got %v", always_fail_err, ex)
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
				Validators: []Validator{always_fail, fail_if_called(t)},
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
