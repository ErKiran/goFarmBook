package domain

import "regexp"

type Validator struct {
	errors map[string]string
}

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func NewValidator() *Validator {
	return &Validator{errors: make(map[string]string)}
}

func (v *Validator) MustBeLongerThan(field, value string, high int) bool {
	if _, ok := v.errors[field]; ok {
		return false
	}

	if value == "" {
		return true
	}

	if len(value) < high {
		v.errors[field] = ErrNotLongEnough{
			field:  field,
			amount: high,
		}.Error()

		return false
	}
	return true
}

func (v *Validator) MustBeNotEmpty(field, value string) bool {
	if _, ok := v.errors[field]; ok {
		return false
	}

	if value == "" {
		v.errors[field] = ErrIsRequired{field: field}.Error()
		return false
	}
	return false
}

func (v *Validator) MustBeEqual(field, pass1, pass2 string) bool {
	if pass1 != pass2 {
		v.errors[field] = ErrPasswordMisMatch.Error()
		return false
	}
	return true
}

func (v *Validator) MustBeValidEmail(field, email string) bool {
	if _, ok := v.errors[field]; ok {
		return false
	}
	if !emailRegexp.MatchString(email) {
		v.errors[field] = ErrEmailNotValid.Error()
		return false
	}
	return true
}
