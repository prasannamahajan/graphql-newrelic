package schema

type Car struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	Colour  string `json:"colour"`
	Abs     bool   `json:"abs"`
	Price   int    `json:"price"`
}

type User struct {
	Name     string `json:"name"`
	MobileNo string `json:"mobileno"`
}

type Query struct {
	Car  Car  `json:"car"`
	User User `jsob:"user"`
}

type QueryCarArgs struct {
	Name string `json:"name"`
}

func (q Query) ArgsForCar() QueryCarArgs {
	return QueryCarArgs{}
}

type CarPriceArgs struct {
	Unit *string `json:"unit"`
}

func (c Car) ArgsForPrice() CarPriceArgs {
	return CarPriceArgs{}
}
