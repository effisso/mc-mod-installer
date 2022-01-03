package input

import (
	"bufio"
	"fmt"
	"io"
)

type Prompt interface {
	GetInput(w io.Writer, r io.Reader) (string, error)
}

type prompt struct {
	PromptText  string
	Validators  []Validator
	CaptureLine bool
}

func NewLinePrompt(promptText string, validators ...Validator) Prompt {
	return prompt{
		PromptText:  promptText,
		Validators:  validators,
		CaptureLine: true,
	}
}

func NewYesNoPrompt(promptText string) Prompt {
	return &prompt{
		PromptText:  promptText,
		Validators:  []Validator{NewRegexValidator("^[yYnN]$", "Expected y/n")},
		CaptureLine: false,
	}
}

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
