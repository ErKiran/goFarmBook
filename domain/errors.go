package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNoResult                     = errors.New("no result")
	ErrUserWithEmailAlreadyExist    = errors.New("user with email already exist")
	ErrUserWithUserNameAlreadyExist = errors.New("user with username already exist")
	ErrEmailNotValid                = errors.New("Email isn't valid email address")
	ErrPasswordMisMatch             = errors.New("Password and Confirm Password doesn't match")
)

type ErrNotLongEnough struct {
	field  string
	amount int
}

func (e ErrNotLongEnough) Error() string {
	return fmt.Sprintf("%v not long emough, %d characters is required", e.field, e.amount)
}

type ErrIsRequired struct {
	field string
}

func (e ErrIsRequired) Error() string {
	return fmt.Sprintf("%v is required", e.field)
}

func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}
