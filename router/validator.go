package router

import "github.com/go-playground/validator/v10"

// NewValidator: 유효성 검증을 위한 객체 생성
func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

type Validator struct {
	validator *validator.Validate
}

// Validate: Validation을 수행하고 그 결과를 반환
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
