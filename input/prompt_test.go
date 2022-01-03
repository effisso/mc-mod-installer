package input_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"mcmods/input"
	. "mcmods/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Prompts", func() {
	var inBuffer *bytes.Buffer
	var outBuffer *bytes.Buffer

	BeforeEach(func() {
		InitTestData()
		outBuffer = bytes.NewBufferString("")
		inBuffer = bytes.NewBufferString("")
	})

	Describe("Line Prompt", func() {
		promptText := "enter something: "

		It("should read from and write to the given buffers", func() {
			inputText := "777"
			inBuffer.WriteString(inputText + "\n")
			p := input.NewLinePrompt(promptText, &input.NoOpValidator{})

			str, err := p.GetInput(outBuffer, inBuffer)

			Expect(err).To(BeNil())
			Expect(str).To(Equal(inputText))
			verifyOutput(outBuffer, promptText)
		})

		It("validates the exact user input", func() {
			inputText := "user input 1"
			validator := &spyValidator{
				ExpectedInput: inputText,
			}
			inBuffer.Write([]byte(inputText + "\n"))
			p := input.NewLinePrompt(promptText, validator)

			str, err := p.GetInput(outBuffer, inBuffer)

			Expect(err).To(BeNil())
			Expect(str).To(Equal(inputText))
			Expect(validator.Visited).To(BeTrue())
		})

		It("prompts multiple times for failed validations", func() {
			validator := &countingValidator{FailureCount: 2}
			inputText1 := "user input 1 (fail)"
			inputText2 := "user input 2 (fail)"
			inputText3 := "user input 3 (succeed)"
			inBuffer.Write([]byte(inputText1 + "\n" + inputText2 + "\n" + inputText3 + "\n"))
			p := input.NewLinePrompt(promptText, validator)
			expectedOutput :=
				promptText + "failure 1 of 2\n" +
					promptText + "failure 2 of 2\n" +
					promptText

			str, err := p.GetInput(outBuffer, inBuffer)

			Expect(err).To(BeNil())
			Expect(str).To(Equal(inputText3))
			verifyOutput(outBuffer, expectedOutput)
		})

		It("returns non-validation errors", func() {
			expectedErr := errors.New("non-validation problem")
			validator := &fakeValidator{Return: expectedErr}
			p := input.NewLinePrompt(promptText, validator)

			_, err := p.GetInput(outBuffer, inBuffer)

			Expect(err).To(Equal(expectedErr))
		})

		It("checks input against all the validators", func() {
			// The caveat here is that it only calls a validator if the previous ones were successful
			// So the flow is like this:
			// input 1 -> rejected by validator 1
			// input 2 -> approved by validator 1, rejected by validator 2
			// input 3 -> approved by both validators
			validator1 := &countingValidator{FailureCount: 1}
			validator2 := &countingValidator{FailureCount: 1}
			inputText1 := "user input 1 (fail)"
			inputText2 := "user input 2 (fail)"
			inputText3 := "user input 3 (succeed)"
			inBuffer.Write([]byte(inputText1 + "\n" + inputText2 + "\n" + inputText3 + "\n"))
			p := input.NewLinePrompt(promptText, validator1, validator2)
			expectedOutput :=
				promptText + "failure 1 of 1\n" +
					promptText + "failure 1 of 1\n" +
					promptText

			str, err := p.GetInput(outBuffer, inBuffer)

			Expect(err).To(BeNil())
			Expect(str).To(Equal(inputText3))
			verifyOutput(outBuffer, expectedOutput)
		})
	})

	Describe("Y/N Prompt", func() {
		promtText := "y/n?"
		p := input.NewYesNoPrompt(promtText)

		It("accepts y/Y/n/N", func() {
			validInputs := []string{"y", "Y", "n", "N"}
			for _, i := range validInputs {
				inBuffer.WriteString(i)

				str, err := p.GetInput(outBuffer, inBuffer)

				Expect(err).To(BeNil())
				Expect(str).To(Equal(i))
			}
		})

		It("ignores characters after the first", func() {
			inBuffer.WriteString("yes")

			str, err := p.GetInput(outBuffer, inBuffer)

			Expect(err).To(BeNil())
			Expect(str).To(Equal("y"))
		})

		It("rejects non y/n chars", func() {
			errMsg := "Expected y/n\n"
			inBuffer.WriteString("h=8*y") // the last char must be valid to exit the prompt loop
			expectedOutput := promtText + errMsg +
				promtText + errMsg +
				promtText + errMsg +
				promtText + errMsg +
				promtText

			str, err := p.GetInput(outBuffer, inBuffer)

			Expect(err).To(BeNil())
			Expect(str).To(Equal("y"))
			verifyOutput(outBuffer, expectedOutput)
		})
	})
})

func verifyOutput(buffer *bytes.Buffer, expected string) {
	out, err := ioutil.ReadAll(buffer)
	Expect(err).To(BeNil())
	strOut := string(out)
	Expect(strOut).To(Equal(expected))
}

type spyValidator struct {
	ExpectedInput string
	Visited       bool
	Return        error
}

func (v *spyValidator) Validate(input string) error {
	v.Visited = true
	Expect(input).To(Equal(v.ExpectedInput))
	return v.Return
}

type countingValidator struct {
	FailureCount uint
	currentCount uint
}

func (v *countingValidator) Validate(str string) error {
	v.currentCount++
	if v.currentCount > v.FailureCount {
		return nil
	}
	return &input.ValidationError{Message: fmt.Sprintf("failure %d of %d", v.currentCount, v.FailureCount)}
}

type fakeValidator struct {
	Return error
}

func (v *fakeValidator) Validate(str string) error {
	return v.Return
}
