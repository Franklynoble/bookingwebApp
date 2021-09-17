package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strings"
)

//Form create a custom error message struct, eme
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

//New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		Values: data,
		Errors: map[string][]string{},
	}
}

//  Required checks for the required field
func (f *Form) Required(field ...string) {
	for _, field := range field {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field can not be blank")
		}
	}
}

//MinLength checks for  minLength
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d Characters long", length))
		return false
	}
	return true
}

//Has checks if form field is post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "this field can not be blank")
		return false
	}
	return true
}

// IsEmail checks for email validator
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid Email Address")
	}
}
