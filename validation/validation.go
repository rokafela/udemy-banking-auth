package validation

import (
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/rokafela/udemy-banking-auth/logger"
)

var validate *validator.Validate

func ValidateStruct(target_validate interface{}) interface{} {
	english := en.New()
	universal_translator := ut.New(english, english)
	translator, found := universal_translator.GetTranslator("en")
	if !found {
		logger.Fatal("translator not found")
	}
	validate = validator.New()
	register_err := en_translations.RegisterDefaultTranslations(validate, translator)
	if register_err != nil {
		logger.Fatal(register_err.Error())
	}
	_ = validate.RegisterTranslation("required", translator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
	validation_error := validate.Struct(target_validate)
	if validation_error != nil {
		if _, ok := validation_error.(*validator.InvalidValidationError); ok {
			fmt.Println(validation_error)
			return validation_error
		}

		for _, err := range validation_error.(validator.ValidationErrors) {
			// fmt.Println(err.Namespace())
			// fmt.Println(err.Field())
			// fmt.Println(err.StructNamespace())
			// fmt.Println(err.StructField())
			// fmt.Println(err.Tag())
			// fmt.Println(err.ActualTag())
			// fmt.Println(err.Kind())
			// fmt.Println(err.Type())
			// fmt.Println(err.Value())
			// fmt.Println(err.Param())
			fmt.Println(err)
		}

		return validation_error
	}

	return nil
}
