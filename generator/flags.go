package generator

import (
	"github.com/brianvoe/gofakeit/v7"
)

func (g *Generator) RegisterFlags() {
	gen := make(Gen)

	gen["randUuid"] = func() any {
		return gofakeit.UUID()
	}

	gen["randString"] = func() any {
		return gofakeit.Word()
	}

	gen["randNum"] = func() any {
		return gofakeit.Int()
	}

	gen["randBool"] = func() any {
		return gofakeit.Bool()
	}

	// будущая дата
	gen["randFDate"] = func() any {
		return gofakeit.FutureDate()
	}

	// прошедшая дата
	gen["randPDate"] = func() any {
		return gofakeit.PastDate()
	}

	g.flags = gen
}
