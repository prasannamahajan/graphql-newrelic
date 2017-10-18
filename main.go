package main

import (
	"encoding/json"
	//"fmt"
	"github.com/arvitaly/go-graphql-tools"
	"github.com/arvitaly/graphql"
	"log"
	"os"
)

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

func GetCar(name string) Car {
	return Car{Name: name, Company: "bmw", Colour: "Black", Abs: true, Price: 500000}
	//	return data[name]
}

func GetUser() User {
	return User{
		Name:     "Prasanna",
		MobileNo: "8800220011",
	}
}

func NewRouter() *tools.Router {
	router := tools.NewRouter()
	router.Query("Query.Car", func() (Car, error) {
		return GetCar("b500"), nil
	})
	router.Query("Query.User", func() (User, error) {
		return GetUser(), nil
	})
	return router
}

func main() {
	router := NewRouter()
	gen := tools.NewGenerator(router)
	query := gen.GenerateObject(Query{})
	//	mutation := gen.GenerateObject(Mutation{})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: query,
		//		Mutation: mutation,
	})
	if err != nil {
		log.Fatal(err)
	}
	q := `query { 
			car{
				name
				company
			}
			user{
				name
			}
		}`
	res := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: q,
	})
	if res.HasErrors() {
		log.Fatalf("Result has errors: %v", res.Errors)
	}
	json.NewEncoder(os.Stdout).Encode(res)
}
