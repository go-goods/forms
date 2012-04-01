package forms

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

type LoaderFunc func(map[string]interface{}) (interface{}, map[string]error, error)

func (l LoaderFunc) Load(in map[string]interface{}) (out interface{}, errs map[string]error, err error) {
	out, errs, err = l(in)
	return
}
