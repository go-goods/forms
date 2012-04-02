package forms

import "net/http"

//Result represents the data returned by loading a form on a request.
type Result struct {
	Value  interface{}
	Values map[string]string
	Errors map[string]error
	Err    error
}

//Form is a type that represents a form that can load requests and produce results
//based on the fields and loader.
type Form struct {
	Fields []Field
	Loader Loader
}

//Load takes a request and returns a result representing the fields having been
//validated and converted. If no errors occur during validation/conversion of
//every field, then the Loader is used to do final validation and conversion
//into a value suitable for processing on.
func (f *Form) Load(req *http.Request) (r *Result) {
	//set up map for errors, values in from the form, and converted values
	r = &Result{
		Values: map[string]string{},
		Errors: map[string]error{},
	}
	conv := map[string]interface{}{}

fields:
	for _, fi := range f.Fields {
		//grab the value for the field and store it
		in := req.FormValue(fi.Name)
		r.Values[fi.Name] = in

		//run all the validators on the value
		val, err := fi.Validate(in)
		if err != nil {
			r.Errors[fi.Name] = err
			continue fields
		}

		//convert the value if we have a converter and store the value
		if fi.Converter != nil {
			c, err := fi.Converter.Convert(val)
			if err != nil {
				r.Errors[fi.Name] = err
				continue fields
			}
			conv[fi.Name] = c
		} else {
			conv[fi.Name] = val
		}
	}

	//if we got any errors dont try to load in a value
	if len(r.Errors) > 0 {
		return
	}

	//If we have a loader, run it, otherwise just return the map of converted
	//values
	if f.Loader != nil {
		r.Value, r.Errors, r.Err = f.Loader.Load(conv)
	} else {
		r.Value = conv
	}

	return
}

//Field represents an individual field in a form. It has a list of Validators
//that will be called in order, exiting early when one of them returns a non-nil
//error. If no errors occured and a Converter is set, then the Convert method is
//called on the final value.
type Field struct {
	Name       string
	Validators []Validator
	Converter  Converter
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

//Validator is a type that can validate a value and return any validation errors
//and a new value.
type Validator interface {
	Validate(in string) (out string, err error)
}

//Converter is a type that converts a string value into something else, returning
//any errors in the conversion.
type Converter interface {
	Convert(in string) (out interface{}, err error)
}

//Loader takes a map of values and returns a value that represents that map, any
//errors constructing that value on a per field name basis, and an error that
//doesn't correspond to any specific field.
type Loader interface {
	Load(in map[string]interface{}) (value interface{}, errs map[string]error, err error)
}
