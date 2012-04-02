package forms

//ValidatorFunc is a function that conforms to the Validator interface. Similar
//to the http.Handler/http.HandlerFunc in net/http
type ValidatorFunc func(string) (string, error)

//Validate implements the Validator interface for ValidatorFunc
func (v ValidatorFunc) Validate(in string) (out string, err error) {
	out, err = v(in)
	return
}

//ConverterFunc is a function that conforms to the Converter interface. Similar
//to the http.Handler/http.HandlerFunc in net/http
type ConverterFunc func(string) (interface{}, error)

//Convert implements the Converter interface for ConverterFunc
func (c ConverterFunc) Convert(in string) (out interface{}, err error) {
	out, err = c(in)
	return
}

//LoaderFunc is a function that conforms to the Loader interface. Similar
//to the http.Handler/http.HandlerFunc in net/http
type LoaderFunc func(map[string]interface{}) (interface{}, map[string]error, error)

//Load implements the Loader interface for LoaderFunc
func (l LoaderFunc) Load(in map[string]interface{}) (out interface{}, errs map[string]error, err error) {
	out, errs, err = l(in)
	return
}

//Assert that our Funcs are of the correct interface type
var (
	_ Loader    = LoaderFunc(nil)
	_ Converter = ConverterFunc(nil)
	_ Validator = ValidatorFunc(nil)
)
