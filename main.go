package main

import (
	"fmt"
	"github.com/arvitaly/go-graphql-tools"
	"github.com/arvitaly/graphql"
	"log"
)

func main() {
	router := NewRouter()
	gen := tools.NewGenerator(router)
	query := gen.GenerateObject(Query{})
	mutation := gen.GenerateObject(Mutation{})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	})
	if err != nil {
		log.Fatal(err)
	}
	q := "test"
	res := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: q,
	})
	if res.HasErrors() {
		log.Fatalf("Result has errors: %v", res.Errors)
	}
	fmt.Println(res)
}
