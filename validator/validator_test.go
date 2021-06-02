package validator_test

import (
	"testing"

	"github.com/go-chassis/foundation/validator"
	"github.com/stretchr/testify/assert"
)

type student struct {
	Name    string `validate:"kieTest"`
	Address string `validate:"alpha,min=2,max=4"`
}

func TestNewValidator(t *testing.T) {
	r := validator.NewRegexRule("kieTest", `^[a-zA-Z0-9]*$`)
	valid := validator.NewValidator()
	err := valid.RegisterRule(r)
	assert.Nil(t, err)
	assert.Nil(t, valid.AddErrorTranslation4Tag("min"))
	assert.Nil(t, valid.AddErrorTranslation4Tag("max"))

	s := &student{Name: "a1", Address: "abc"}
	err = valid.Validate(s)
	assert.Nil(t, err)

	s = &student{Name: "a1-", Address: "abc"}
	err = valid.Validate(s)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "field: student.Name, rule: ^[a-zA-Z0-9]*$")

	s = &student{Name: "a1", Address: "abcde"}
	err = valid.Validate(s)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "field: student.Address, rule: max = 4")
}
