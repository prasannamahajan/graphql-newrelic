package main

import (
	"encoding/json"
	//"fmt"
	"fmt"
	"github.com/arvitaly/go-graphql-tools"
	"github.com/arvitaly/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"graphql-newrelic/data"
	"graphql-newrelic/relicconf"
	"graphql-newrelic/router"
	"graphql-newrelic/schema"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func InitSchema() graphql.Schema {
	router := router.NewRouter()
	gen := tools.NewGenerator(router)
	query := gen.GenerateObject(schema.Query{})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: query,
	})
	if err != nil {
		log.Fatal(err)
	}
	return schema
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	fmt.Println("query is : ", query)
	result := graphql.Do(graphql.Params{
		RequestString: query,
		Schema:        schema,
	},
	)
	if len(result.Errors) > 0 {
		fmt.Println("some error occured")
	}
	return result
}

func getTypesFromQuery(q string) (string, error) {
	source := source.NewSource(&source.Source{
		Body: []byte(q),
		Name: "GraphQL request",
	})
	querytypes := make([]string, 0)
	AST, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	for _, definition := range AST.Definitions {
		switch definition := definition.(type) {
		case *ast.OperationDefinition:
			for _, selection := range definition.SelectionSet.Selections {
				switch selection := selection.(type) {
				case *ast.Field:
					querytypes = append(querytypes, selection.Name.Value)
				}
			}
		}
	}
	return strings.Join(querytypes, ","), nil
}

var gquery string

func queryHandler(w http.ResponseWriter, r *http.Request) {
	result := executeQuery(gquery, Schema)
	json.NewEncoder(w).Encode(result)
}

func slog(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		gquery = string(body)
		types, _ := getTypesFromQuery(gquery)
		fmt.Println("types :-", types)
		app := relicconf.GetRelicApp()
		txn := app.StartTransaction(types, nil, nil)
		defer txn.End()
		h.ServeHTTP(w, r) // call original
	})
}

var Schema graphql.Schema

func dummyQuery() {
	q := `query {
				user{
					name
				}
			}`
	res := executeQuery(q, Schema)
	json.NewEncoder(os.Stdout).Encode(res)
}
func init() {
	Schema = InitSchema()
	relicconf.InitNewRelic()
	data.ImportData()
	//	dummyQuery()
}

func main() {
	http.HandleFunc("/graphql", slog(queryHandler))
	fmt.Println("Starting Server")
	http.ListenAndServe(":8080", nil)
}
