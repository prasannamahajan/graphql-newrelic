package data

import "graphql-newrelic/schema"

var Data map[string]schema.Car

func GetCar(name string) schema.Car {
	return Data[name]
}

func GetUser() schema.User {
	return schema.User{
		Name:     "Prasanna",
		MobileNo: "8800220011",
	}
}

func ImportData() {
	Data = make(map[string]schema.Car)
	Data["b500"] = schema.Car{Name: "b500", Company: "bmw", Colour: "blue", Abs: true, Price: 3000000}
	Data["indigo"] = schema.Car{Name: "indigo", Company: "tata", Colour: "black", Abs: false, Price: 1000000}
	Data["swift"] = schema.Car{Name: "swift", Company: "maruti", Colour: "grey", Abs: true, Price: 500000}
}
