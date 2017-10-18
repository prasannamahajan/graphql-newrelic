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
	q := `query Q1{ 
			rebels{
				id 
				name
				ships{
					edges{
						node{
							id
							name
						}
					}
				}
			} 
			empire{
				id 
				name
			}				 
			car{
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
