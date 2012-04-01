package forms

type Form struct {
	Fields []Field
	Loader Loader
}

type Field struct {
	Name, Display string
	Validators    []Validator
	Converter     Converter
}

type Validator interface {
	Validate(string) (string, error)
}

type Converter interface {
	Convert(string) (interface{}, error)
}

type Loader interface {
	Load(map[string]interface{}) (interface{}, error, map[string]error)
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

type LoaderFunc func(map[string]interface{}) (interface{}, error, map[string]error)

func (l LoaderFunc) Load(in map[string]interface{}) (out interface{}, err error, errs map[string]error) {
	out, err, errs = l(in)
	return
}
