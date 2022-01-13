package api

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type beerValidator struct {
	*validator.Validate
	t ut.Translator
}

func (v *beerValidator) translate(err error) []string {
	var result []string
	fmtErrs := err.(validator.ValidationErrors).Translate(v.t)

	for _, e := range fmtErrs {
		result = append(result, e)
	}

	return result
}

func newValidator() *beerValidator {
	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")

	v := validator.New()
	en_translations.RegisterDefaultTranslations(v, trans)

	v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} must have a value", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	return &beerValidator{v, trans}
}
