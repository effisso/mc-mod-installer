package input

import (
	"bufio"
	"fmt"
	"io"
)

// Prompt is a way of soliciting input from the user by posing a question and
// reading input
type Prompt interface {
	// GetInput prints to the user on the writer, and reads their input from the reader
	GetInput(w io.Writer, r io.Reader) (string, error)
}

type prompt struct {
	PromptText  string
	Validators  []Validator
	CaptureLine bool
}

// NewLinePrompt creates a new Prompt which reads all user input until a
// carriage return
func NewLinePrompt(promptText string, validators ...Validator) Prompt {
	return prompt{
		PromptText:  promptText,
		Validators:  validators,
		CaptureLine: true,
	}
}

// NewYesNoPrompt creates a new Prompt which only accepts Y/y/N/n for a response
func NewYesNoPrompt(promptText string) Prompt {
	return &prompt{
		PromptText:  promptText,
		Validators:  []Validator{NewRegexValidator("^[yYnN]$", "Expected y/n")},
		CaptureLine: false,
	}
}

// GetInput prints to the user on the writer, and reads their input from the reader
func (p prompt) GetInput(w io.Writer, r io.Reader) (string, error) {
	doPrompt := true
	response := ""
	var err error
	scanner := bufio.NewScanner(r)

	for doPrompt {
		fmt.Fprint(w, p.PromptText)

		if p.CaptureLine {
			scanner.Scan()
			response = scanner.Text()
		} else {
			var c byte
			fmt.Fscanf(r, "%c", &c)
			response = string(c)
		}

		for _, validator := range p.Validators {
			err = validator.Validate(response)

			if err == nil {
				continue
			} else if _, ok := err.(*ValidationError); ok {
				fmt.Fprintln(w, err.Error())
				break
			} else {
				return "", err // non-recoverable error, just return it
			}
		}

		if err == nil {
			doPrompt = false
		}
	}

	return response, err
}
