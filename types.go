package forms

import (
	"fmt"
	"net/http"
)

type Result struct {
	Value  interface{}
	Values map[string]string
	Errors map[string]error
}

type Form struct {
	Fields []Field
	Loader Loader
}

func (f *Form) Load(req *http.Request) (r Result, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("load: %v", e)
		}
	}()

	r = Result{
		Values: map[string]string{},
		Errors: map[string]error{},
	}

	//set up map for errors, values in from the form, and converted values
	conv := map[string]interface{}{}

	for _, fi := range f.Fields {
		//grab the value for the field and store it
		in := req.FormValue(fi.Name)
		r.Values[fi.Name] = in

		//run all the validators on the value
		val, err := fi.Validate(in)
		if err != nil {
			r.Errors[fi.Name] = err
			continue
		}

		//convert the value
		c, err := fi.Converter.Convert(val)
		if err != nil {
			r.Errors[fi.Name] = err
			continue
		}
		//store the converted value
		conv[fi.Name] = c
	}

	//if we got any errors dont try to load in a value
	if len(r.Errors) > 0 {
		return
	}

	//load in the value and perform any final validations
	r.Value, r.Errors = f.Loader.Load(conv)
	return
}

type Field struct {
	Name, Display string
	Validators    []Validator
	Converter     Converter
}

//Validate runs all of the validators on a value, exiting if any of the
//validators return an error.
func (f *Field) Validate(in string) (final string, err error) {
	final = in
	for _, v := range f.Validators {
		final, err = v.Validate(final)
		if err != nil {
			return
		}
	}
	return
}

type Validator interface {
	Validate(in string) (string, error)
}

type Converter interface {
	Convert(in string) (interface{}, error)
}

type Loader interface {
	Load(in map[string]interface{}) (interface{}, map[string]error)
}

type ValidatorFunc func(string) (string, error)

func (v ValidatorFunc) Validate(in string) (out string, err error) {
	out, err = v(in)
	return
}

type ConverterFunc func(string) (interface{}, error)

func (c ConverterFunc) Convert(in string) (out interface{}, err error) {
	out, err = c(in)
	return
}

type LoaderFunc func(map[string]interface{}) (interface{}, map[string]error)

func (l LoaderFunc) Load(in map[string]interface{}) (out interface{}, errs map[string]error) {
	out, errs = l(in)
	return
}
