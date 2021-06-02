package validator

import (
	"github.com/go-playground/validator"
	"regexp"
)

// RegexValidateRule contains an validate tag's info
type RegexValidateRule struct {
	tag   string
	regex *regexp.Regexp
}

// Validate validates string
func (r *RegexValidateRule) Validate(s string) bool {
	return r.regex.MatchString(s)
}

func (r *RegexValidateRule) validateFL(fl validator.FieldLevel) bool {
	return r.Validate(fl.Field().String())
}

// Tag returns the validate rule's tag
func (r *RegexValidateRule) Tag() string {
	return r.tag
}

// Explain explains the rule
func (r *RegexValidateRule) Explain() string {
	explain := r.regex.String()
	return explain
}

// NewRegexRule news a rule
func NewRegexRule(tag, regexStr string) *RegexValidateRule {
	return &RegexValidateRule{
		tag:   tag,
		regex: regexp.MustCompile(regexStr),
	}
}
