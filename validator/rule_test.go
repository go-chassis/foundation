package validator_test

import (
	"testing"

	"github.com/go-chassis/foundation/validator"
	"github.com/stretchr/testify/assert"
)

func TestNewRule(t *testing.T) {
	rule := validator.NewRegexRule("t", `^[a-zA-Z0-9]*$`)
	assert.Equal(t, "t", rule.Tag())
	assert.Equal(t, `^[a-zA-Z0-9]*$`, rule.Explain())
	assert.True(t, rule.Validate("ab"))
	assert.True(t, rule.Validate(""))
	assert.True(t, rule.Validate("a"))
	assert.True(t, rule.Validate("abcde"))
	assert.True(t, rule.Validate("abcdefg12345678"))
	assert.False(t, rule.Validate("ab-"))
}
