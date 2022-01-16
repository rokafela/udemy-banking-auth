package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/rokafela/udemy-banking-auth/logger"
)

// initialize validator in package global scope
var validate *validator.Validate

// ValidateStruct attempts to validate any struct using their validate tag
// target_validate uses interface type to support any struct
// the return uses interface to support nil or slice of strings
func ValidateStruct(target_validate interface{}) interface{} {
	// initialize translator
	english := en.New()
	universal_translator := ut.New(english, english)
	translator, found := universal_translator.GetTranslator("en")
	if !found {
		logger.Fatal("translator not found")
	}

	// initialize validator
	validate = validator.New()

	// set validator to use json tag instead of struct property name
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// register translator
	register_err := en_translations.RegisterDefaultTranslations(validate, translator)
	if register_err != nil {
		logger.Fatal(register_err.Error())
	}

	// // register translation
	// _ = validate.RegisterTranslation("required", translator, func(ut ut.Translator) error {
	// 	return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	// }, func(ut ut.Translator, fe validator.FieldError) string {
	// 	t, _ := ut.T("required", fe.Field())
	// 	return t
	// })

	// validate the struct
	validation_error := validate.Struct(target_validate)

	// initialize validation message
	var validation_err_message []string

	// check validation error
	if validation_error != nil {
		// compile validation error message
		for _, err := range validation_error.(validator.ValidationErrors) {
			validation_err_message = append(validation_err_message, err.Translate(translator))
		}
		return validation_err_message
	}

	// no error
	return nil
}
