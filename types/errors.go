package types

import (
	"fmt"
)

type TemplateNotFoundErr struct {
	template string
}

func NewTemplateNotFoundErr(t string) TemplateNotFoundErr { return TemplateNotFoundErr{template: t} }

func (e TemplateNotFoundErr) Error() string {
	return fmt.Sprintf("template '%s' not found", e.template)
}
