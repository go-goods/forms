package forms

import "net/http"

type Result struct {
	Value  interface{}
	Values map[string]string
	Errors map[string]error
	Err    error
}

type Form struct {
	Fields []Field
	Loader Loader
}

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
	Load(in map[string]interface{}) (interface{}, map[string]error, error)
}
