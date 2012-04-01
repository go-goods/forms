package forms

import (
	"net/url"
	"testing"
)

var (
	int_converter ConverterFunc = func(in string) (out interface{}, err error) {
		out = 2
		return
	}
)

func TestConvertInt(t *testing.T) {
	f := &Form{
		Fields: []Field{
			{Name: "foo", Converter: int_converter},
		},
	}
	res := f.Load(create_req(url.Values{
		"foo": {"bar"},
	}))
	if ex, ok := res.Errors["foo"]; ok || ex != nil {
		t.Fatalf("Expected %v. Got %v", nil, ex)
	}
	rval := res.Value.(map[string]interface{})
	if ex, ok := rval["foo"]; !ok || ex.(int) != 2 {
		t.Fatalf("Expected %v. Got %v", 2, ex)
	}
}
