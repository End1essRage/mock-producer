package generator

import (
	"github.com/brianvoe/gofakeit/v7"
)

type Gen map[string]func() any

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
		return gofakeit.FutureDate().Format("2006-01-02T15:04:05")
	}

	// прошедшая дата
	gen["randPDate"] = func() any {
		return gofakeit.PastDate().Format("2006-01-02T15:04:05")
	}

	statuses := []string{"Черновик", "На проверку", "Назначен проверяющий", "В работе", "Проверен СБ",
		"Черный список", "Отложено", "На проверку повторно"}
	gen["typeStatus"] = func() any {
		index := gofakeit.IntRange(0, len(statuses))
		return statuses[index]
	}

	ownerShips := []string{"ООО", "АО", "ИП", "ЧП", "Физическое лицо", "Иностранное юридическое лицо",
		"Иностранное физическое лицо", "МУП", "ПИФ", "ЗПИФ",
		"Некоммерческая организация", "Другое", "ПАО", "ОАО"}
	gen["typeOwnerShip"] = func() any {
		index := gofakeit.IntRange(0, len(ownerShips))
		return ownerShips[index]
	}

	taxModes := []string{"Классическая система н/о", "УСН", "ЕНВД", "ЕСХН", "Патент"}
	gen["typeTaxMode"] = func() any {
		index := gofakeit.IntRange(0, len(taxModes))
		return taxModes[index]
	}

	/*
		Возможные значения: (это для поля status)
		typeOwnerShip
		Возможные значения:  (это для поля ownership_type)
		typeTaxMode
		Возможные значения: Классическая система н/о / УСН / ЕНВД / ЕСХН / Патент (это для поля tax_mode)
	*/

	g.flags = gen
}
