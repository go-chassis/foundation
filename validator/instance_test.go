package validator_test

import (
	"github.com/go-chassis/foundation/validator"
	"github.com/stretchr/testify/assert"
	"testing"
)

const commonNameRegexString = `^[a-zA-Z0-9]*$|^[a-zA-Z0-9][a-zA-Z0-9_\-.]*[a-zA-Z0-9]$`

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"input ABC should pass", struct {
			A string `validate:"min=1,max=3,commonName"`
		}{A: "ABC"}, false},
		{"input ABCD should return err", struct {
			A string `validate:"min=1,max=3,commonName"`
		}{A: "ABCD"}, true},
		{"input A_C should should pass", struct {
			A string `validate:"min=1,max=3,commonName"`
		}{A: "A_C"}, false},
		{"input empty should return err", struct {
			A string `validate:"min=1,max=3,commonName"`
		}{A: ""}, true},
	}

	assert.NoError(t, validator.RegisterRegexRules([]*validator.RegexValidateRule{
		validator.NewRegexRule("commonName", commonNameRegexString),
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validator.Validate(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
