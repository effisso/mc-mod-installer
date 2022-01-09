package input_test

import (
	"fmt"
	"mcmods/cmd"
	"mcmods/input"
	"mcmods/mc"
	. "mcmods/testdata"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInputs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Input Suite")
}

var _ = Describe("Validations", func() {
	BeforeEach(func() {
		InitTestData()
	})

	Describe("Regex Validator", func() {
		errMsg := "regex validation error"

		When("valid regex", func() {
			var validator *input.RegexValidator

			BeforeEach(func() {
				v := input.NewRegexValidator("^[a-z]+$", errMsg)
				validator, _ = v.(*input.RegexValidator)
			})

			It("should construct a new regex", func() {
				Expect(validator.Regex).ToNot(BeNil())
			})

			It("should return no errors for valid strings", func() {
				Expect(validator.Validate("str")).To(BeNil())
			})

			It("should return the validator's error message for invalid strings", func() {
				err := validator.Validate("A")
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal(errMsg))
			})
		})

		When("invalid regex", func() {
			It("should panic", func() {
				defer func() {
					if r := recover(); r == nil {
						Fail("Did not panic")
					}
				}()

				input.NewRegexValidator("([a-z]+", errMsg)

				Fail("Did not panic")
			})
		})
	})

	Describe("URL Validator", func() {
		var validator *input.UrlValidator
		validUrls := []string{
			"https://www.curseforge.com/minecraft/mc-mods/ducts/download/3571121/file",
			"https://www.curseforge.com/minecraft/mc-mods/ducts",
		}
		invalidUrls := []string{
			"", "hello", ":&&6//wat+2=5%;",
		}

		BeforeEach(func() {
			validator = &input.UrlValidator{}
		})

		It("should return no errors for valid URLs", func() {
			for _, str := range validUrls {
				Expect(validator.Validate(str)).To(BeNil(), fmt.Sprintf("failed to validate: %s", str))
			}
		})

		It("should return errors for invalid URLs", func() {
			var asgnTest *input.ValidationError
			for _, str := range invalidUrls {
				err := validator.Validate(str)
				Expect(err).ToNot(BeNil(), fmt.Sprintf("expected error not returned for: %s", str))
				Expect(err).To(BeAssignableToTypeOf(asgnTest))
			}
		})
	})

	Describe("No-op Validator", func() {
		It("should only return nil", func() {
			validator := &input.NoOpValidator{}
			anything := []string{
				"", " ", "yes", ".dkjv0932 -slfd  \t", "\"",
			}
			for _, str := range anything {
				Expect(validator.Validate(str)).To(BeNil(), fmt.Sprintf("got non-nil error for: %s", str))
			}
		})
	})

	Describe("CLI Name Uniqueness Validator", func() {
		var validator *input.CliNameUniquenessValidator

		BeforeEach(func() {
			validator = &input.CliNameUniquenessValidator{
				GetModMap: func() mc.ModMap { return TestingCliModMap },
			}
		})

		It("should return nil if the name is not in use", func() {
			err := validator.Validate("not-in-use")

			Expect(err).To(BeNil())
		})

		It("should return a validation error if the name is in use", func() {
			var asgnTest *input.ValidationError
			err := validator.Validate(TestingClientMod1.CliName)

			Expect(err).ToNot(BeNil())
			Expect(err).To(BeAssignableToTypeOf(asgnTest))
		})
	})

	Describe("Server Group Name Validator", func() {
		var validator *input.GroupNameValidator

		BeforeEach(func() {
			validator = &input.GroupNameValidator{}
		})

		It("should return nil if the name is valid", func() {
			groups := []string{"required", "optional", "performance", cmd.ServerOnlyGroupKey}

			for _, name := range groups {
				err := validator.Validate(name)

				Expect(err).To(BeNil(), "Validate returned an error for "+name)
			}
		})

		It("should return a validation error if the name is not valid", func() {
			var asgnTest *input.ValidationError
			err := validator.Validate("invalid")

			Expect(err).ToNot(BeNil())
			Expect(err).To(BeAssignableToTypeOf(asgnTest))
		})
	})
})
