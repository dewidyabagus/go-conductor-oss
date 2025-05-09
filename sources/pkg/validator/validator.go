package validator

import (
	"bytes"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	once  sync.Once
	vld   *Validate
	trans ut.Translator
)

type Errors map[string]string

func (e Errors) Error() string {
	buff := bytes.NewBufferString("")

	for _, v := range e {
		_, _ = buff.WriteString(v + "\n")
	}
	return strings.TrimSpace(buff.String())
}

type Validate struct {
	validator *validator.Validate
}

func New(opts ...validator.Option) *Validate {
	once.Do(func() {
		en := en.New()
		uni := ut.New(en, en)
		trans, _ = uni.GetTranslator("en")

		v := validator.New(opts...)
		en_translations.RegisterDefaultTranslations(v, trans)

		vld = &Validate{v}
	})

	return vld
}

func (v *Validate) Struct(val any) error {
	return v.translate(v.validator.Struct(val))
}

func (v *Validate) Var(field any, tag string) error {
	return v.translate(v.validator.Var(field, tag))
}

func (Validate) translate(err error) error {
	if err == nil {
		return nil
	}

	vldErrs, ok := err.(validator.ValidationErrors)
	if !ok || len(vldErrs) == 0 {
		return err
	}

	errs := Errors{}
	for i := 0; i < len(vldErrs); i++ {
		errs[vldErrs[i].Field()] = vldErrs[i].Translate(trans)
	}
	return &errs
}
