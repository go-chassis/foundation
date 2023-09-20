package validator

import (
	"regexp"

	"github.com/go-playground/validator"
)

// RegexValidateRule contains an validate tag's info
type RegexValidateRule struct {
	tag     string
	regex   *regexp.Regexp
	explain string
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
	if r.explain != "" {
		return r.explain
	}
	return r.regex.String()
}

// NewRegexRule news a rule
func NewRegexRule(tag, regexStr string) *RegexValidateRule {
	return &RegexValidateRule{
		tag:   tag,
		regex: regexp.MustCompile(regexStr),
	}
}

// WithExplain customize the explanation
func (r *RegexValidateRule) WithExplain(explain string) *RegexValidateRule {
	r.explain = explain
	return r
}
