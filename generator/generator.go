package generator

import (
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/api"
)

//генерация по предоставленному шаблону
//генерация по template

type Gen map[string]func() any

type Generator struct {
	flags Gen
}

func New() *Generator {
	g := &Generator{}
	g.RegisterFlags()

	return g
}

func (g *Generator) Generate(pattern api.Pattern) api.Pattern {
	return g.processValue(pattern).(api.Pattern)
}

func (g *Generator) processValue(value interface{}) interface{} {
	switch v := value.(type) {
	case api.Pattern:
		return g.processMap(v)
	case map[string]interface{}:
		return g.processMap(v)
	case []interface{}:
		return g.processSlice(v)
	case string:
		if generator, exists := g.flags[v]; exists {
			return generator()
		}
		return v
	default:
		return v
	}
}

func (g *Generator) processMap(m map[string]interface{}) api.Pattern {
	result := make(map[string]interface{})
	for key, value := range m {
		result[key] = g.processValue(value)
	}
	return result
}

func (g *Generator) processSlice(s []interface{}) []interface{} {
	result := make([]interface{}, len(s))
	for i, item := range s {
		result[i] = g.processValue(item)
	}
	return result
}
