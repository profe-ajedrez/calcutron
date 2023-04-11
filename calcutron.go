package calcutron

import (
	"github.com/shopspring/decimal"
)

// DiscountTuple contiene un número que indica la cantidad a descontar
// y un tipo indicando si la cantidad a descontar se descontará como porcentaje  como un monto absoluto
type DiscountTuple struct {
	Guarism decimal.Decimal
	Type    int8
}

// ItemValuable representa una cosa que tiene o puede tener valor unitario, impuestos y descuento
// Generalmente es una abstracción de un elemento del carro de compras
// Puede hacer tantas implementaciones de ItemValuable como necesite para su caso específico.
// Se provee una implementación base en el tipo `ItemValue`
type ItemValuable interface {
	Valuable
	Quantity() decimal.Decimal
}

// Valuable Representa una cosa a la que se le puede dar un valor y calcularle impuestos y descuentos
type Valuable interface {
	UnitValue() decimal.Decimal
	Discount() DiscountTuple
	ChangeRate() decimal.Decimal
	Dependencies() Dependencies

	MakeSubtotalWithoutDiscount() decimal.Decimal
	MakeSubtotal() decimal.Decimal
}

// ItemListValuable representa una lista de `ItemValuable`
type ItemListValuable []ItemValuable


// IDetMaker representa una cosa con la capacidad de calcular subtotales, descuentos e impuestos desde
// los datos base de una línea de detalle.
// Es una interface compuesta a su vez por las interfaces:
//   - Maker
type IDetMaker interface {
	Maker
}

// Maker es una cosa capaz de generar un detalle con los valores subtotales a partir de un ItemValuable
type Maker interface {
	Make(iv ItemValuable) (detailLocal *Detail, detailExp *Detail)
}

// Discounter es una cosa capaz de calcular descuentos sobre un itemValuable
type Discounter interface {
	DiscountGuarism(iv ItemValuable) (decimal.Decimal, decimal.Decimal)
	NetDiscountAmount(iv ItemValuable) decimal.Decimal
	BruteDiscountAmount(iv ItemValuable) decimal.Decimal
	UnitWDiscount(iv ItemValuable) decimal.Decimal
}

// Totalizer es capaz de calcular subtotales a pa rtir de un itemValuable
type Totalizer interface {
	NetTotal(iv ItemValuable) decimal.Decimal
	BruteTotal(iv ItemValuable) decimal.Decimal
	NetWDTotal(iv ItemValuable) decimal.Decimal
	BruteWDTotal(iv ItemValuable) decimal.Decimal
}

// Taxalizer es capaz de calcular impuestos y taxables a partir de un itemValuable
// Taxalizer debe poder diferenciar y calcular como mínimo:
//   - impuestos porcentuales y de monto, estos últimos que se multipliquen por cantidad y que no
//   - impuestos porcentuales y de monto que sean sobreimpuesto,estos últimos que se multipliquen por cantidad y que no
//   - impuesto porcentuales y de monto que no sean sobreimpuesto y no seban usarse
//     para calcular los sobre impuestos (caso impuesto a las bolsas Perú)
type Taxalizer interface {
	TaxTotal(iv ItemValuable) decimal.Decimal
	TaxTotalWithoutDisc(iv ItemValuable) decimal.Decimal
}

// ListMaker es un cosa capaz de generar una lista de detalles a partir de ItemListValuable
type ListMaker interface {
	Make(ivs ItemListValuable) ([]*Detail, []*Detail)
}


// Detail permite almacenar los valores subtotales calculados a partir de una
// implementación de ItemValue
type Detail struct {
	Net           decimal.Decimal `json:"net"`
	Brute         decimal.Decimal `json:"brute"`
	Tax           decimal.Decimal `json:"tax"`
	Discount      decimal.Decimal `json:"netDiscount"`
	DiscountBrute decimal.Decimal `json:"bruteDiscount"`
	NetWd         decimal.Decimal `json:"netWithoutDiscount"`
	BruteWd       decimal.Decimal `json:"bruteWithoutDiscount"`
	TaxWd         decimal.Decimal `json:"taxWithoutDiscount"`
}
