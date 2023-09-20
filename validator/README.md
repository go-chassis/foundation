# Validator

Enhanced the https://github.com/go-playground/validator library, support:
- Declare the customizable tag of regular expressions

#### Usage

First, declare a struct and add **`validate`** tag.
```go
type Struct struct {
    Field `validate:"min=1,max=6"`
}
```

Or declare a regular expression in some complex scenes.
```go
const commonNameRegexString = `^[a-zA-Z0-9][a-zA-Z0-9_\-.]*[a-zA-Z0-9]$`

type Struct struct {
    Field `validate:"min=1,max=6,commonName"`
}

func init() {
    validator.RegisterRegexRules([]*validator.RegexValidateRule{
		validator.NewRegexRule("commonName", commonNameRegexString),
	})
}
```

Finally, add validator check to specific business code.
```go
it := &Struct {
    Field: "AB**CD"
}
if err := validator.Validate(it); err != nil {
    // handle err
}
```