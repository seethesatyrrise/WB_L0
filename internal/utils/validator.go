package utils

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"http-nats-psql/internal/models"
)

type Validator struct {
	V *validator.Validate
}

func NewValidator() (*Validator, error) {
	return &Validator{V: validator.New()}, nil
}

func (v *Validator) Validate(msg []byte) (bool, error) {
	parsed := &models.Order{Items: []models.Items{models.Items{}}}
	json.Unmarshal(msg, parsed)
	err := v.V.Struct(parsed)
	if err != nil {
		return false, err
	}
	return true, nil
}
