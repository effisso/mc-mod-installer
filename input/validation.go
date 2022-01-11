package input

import (
	"mcmods/mc"
	"net/url"
	"regexp"
)

// Validator validates user input
type Validator interface {
	// Validate returns an error if the given input string is not valid
	Validate(input string) error
}

// NoOpValidator does nothing but return a nil error during validation
type NoOpValidator struct{}

// Validate does nothing, just returns nil
func (v NoOpValidator) Validate(input string) error {
	return nil
}

// RegexValidator checks the input against a regular expression to check validity
type RegexValidator struct {
	Regex      regexp.Regexp
	errMessage string
}

// NewRegexValidator creates a new instance of RegexValidator with a compiled
// regular expression for the exp argument
func NewRegexValidator(exp string, errMsg string) Validator {
	return &RegexValidator{
		Regex:      *regexp.MustCompile(exp),
		errMessage: errMsg,
	}
}

// Validate returns an error if the given input string does not match the regex
func (v *RegexValidator) Validate(input string) error {
	if !v.Regex.MatchString(input) {
		return &ValidationError{Message: v.errMessage}
	}
	return nil
}

// URLValidator makes sure that the URL is valid, but doesn't check reachability
type URLValidator struct{}

// Validate returns an error if the given input string is not parsable to a URL
func (v *URLValidator) Validate(input string) error {
	if _, err := url.ParseRequestURI(input); err != nil {
		return &ValidationError{Message: "Invalid URL: " + err.Error()}
	}
	return nil
}

// CliNameUniquenessValidator ensures that the given name for the CLI is not
// already in use
type CliNameUniquenessValidator struct {
	GetModMap func() mc.ModMap
}

// Validate returns an error if the given input string is already a mod CLI name
func (v *CliNameUniquenessValidator) Validate(input string) error {
	if _, exists := v.GetModMap()[input]; exists {
		return &ValidationError{Message: "Name is not unique: " + input}
	}
	return nil
}

// GroupNameValidator ensures that the group name provided by the user is valid
type GroupNameValidator struct{}

// Validate returns an error if the given input string is group name doesn't exist
func (v *GroupNameValidator) Validate(input string) error {
	if _, exists := mc.ServerGroups[input]; !exists {
		return &ValidationError{Message: "Unknown server group: " + input}
	}
	return nil
}

// ValidationError describes non-fatal validation problems with the user's input
type ValidationError struct {
	Message string
}

// Error
func (e *ValidationError) Error() string {
	return e.Message
}
