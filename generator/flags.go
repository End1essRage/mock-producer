package generator

import (
	"github.com/brianvoe/gofakeit/v7"
)

func (g *Generator) RegisterFlags() {
	gen := make(Gen)

	gen["randString"] = func() any {
		return gofakeit.Word()
	}

	gen["randNum"] = func() any {
		return gofakeit.Int()
	}

	g.flags = gen
}
