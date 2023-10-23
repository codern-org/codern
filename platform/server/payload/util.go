package payload

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type FieldsSelector struct {
	fields []string
}

func GetFieldSelector(ctx *fiber.Ctx) *FieldsSelector {
	var selector FieldsSelector
	value := ctx.Query("fields")

	if value == "" {
		return &selector
	}

	values := strings.Split(value, ",")
	for i := 0; i < len(values); i++ {
		selector.fields = append(selector.fields, values[i])
	}
	return &selector
}

func (p *FieldsSelector) Has(field string) bool {
	for i := 0; i < len(p.fields); i++ {
		if p.fields[i] == field {
			return true
		}
	}
	return false
}
