package validator_test

import (
	"testing"

	"github.com/go-chassis/foundation/validator"
	"github.com/stretchr/testify/assert"
)

func TestNewRule(t *testing.T) {
	rule1 := validator.NewRegexRule("t", `^[a-zA-Z0-9]*$`)
	assert.Equal(t, "t", rule1.Tag())
	assert.Equal(t, `^[a-zA-Z0-9]*$`, rule1.Explain())
	assert.True(t, rule1.Validate("ab"))
	assert.True(t, rule1.Validate(""))
	assert.True(t, rule1.Validate("a"))
	assert.True(t, rule1.Validate("abcde"))
	assert.True(t, rule1.Validate("abcdefg12345678"))
	assert.False(t, rule1.Validate("ab-"))

	// test NewRegexRule with explain
	rule2 := validator.NewRegexRule("t", `^[a-zA-Z0-9]*$`).WithExplain("some rule you don't want to expose")
	assert.Equal(t, "some rule you don't want to expose", rule2.Explain())
}
