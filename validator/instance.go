package validator

var GlobalValidator = NewValidator()

func init() {
	if err := Wrap3rdTagsTranslation(); err != nil {
		panic(err)
	}
}

func RegisterRegexRules(rules []*RegexValidateRule) error {
	for _, r := range rules {
		if err := GlobalValidator.RegisterRule(r); err != nil {
			return err
		}
	}
	return nil
}

// Validate validates data
func Validate(v interface{}) error {
	return GlobalValidator.Validate(v)
}
