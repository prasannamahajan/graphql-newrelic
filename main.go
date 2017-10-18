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
	Name   string `json:"name"`
	Wheels int    `json:"wheels"`
}

type Query struct {
	Car Car `json:"car"`
}

func NewRouter() *tools.Router {
	router := tools.NewRouter()
	router.Query("Query.Car", func() (Car, error) {
		return Car{Name: "bmw", Wheels: 4}, nil
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
	q := `query Q2{ 
			car{
				name
				wheels
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
