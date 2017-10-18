package main

import (
	"encoding/json"
	//"fmt"
	"fmt"
	"github.com/arvitaly/go-graphql-tools"
	"github.com/arvitaly/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/printer"
	"github.com/graphql-go/graphql/language/source"
	newrelic "github.com/newrelic/go-agent"
	"io/ioutil"
	"log"
	"net/http"
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

type QueryCarArgs struct {
	Name string `json:"name"`
}

func (q Query) ArgsForCar() QueryCarArgs {
	return QueryCarArgs{}
}

func GetCar(name string) Car {
	return data[name]
}

func GetUser() User {
	return User{
		Name:     "Prasanna",
		MobileNo: "8800220011",
	}
}

func NewRouter() *tools.Router {
	router := tools.NewRouter()
	router.Query("Query.User", func() (User, error) {
		return GetUser(), nil
	})
	router.Query("Query.Car", func(q Query, args QueryCarArgs) (interface{}, error) {
		return GetCar(args.Name), nil
	})
	return router
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

func customProcess(p graphql.Schema, q string) {
	source := source.NewSource(&source.Source{
		Body: []byte(q),
		Name: "GraphQL request",
	})
	AST, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		fmt.Println(err)
	}
	as := printer.Print(AST.Definitions[0])
	fmt.Println(as)
	fmt.Println("--------")
	for _, definition := range AST.Definitions {
		switch definition := definition.(type) {
		case *ast.OperationDefinition:
			fmt.Println("--------")
			fmt.Println("In ast operation definition")
			fmt.Println("Def name", definition.Name)
			fmt.Println("Def kind", definition.Kind)
			fmt.Println("Def directives", definition.GetDirectives())
			fmt.Println("Def Operation", definition.GetOperation())
			//fmt.Println("Def SelectionSet", definition.Se
			if definition.GetName() != nil && definition.GetName().Value != "" {
				fmt.Println("[", definition.GetName().Value, "]")
			}
			fmt.Println("--------")
		case *ast.FragmentDefinition:
			fmt.Println("In FragmentDefinition")
			/*
				key := ""
				if definition.GetName() != nil && definition.GetName().Value != "" {
					key = definition.GetName().Value
				}
						fragments[key] = definition
					default:
						return nil, fmt.Errorf("GraphQL cannot execute a request containing a %v", definition.GetKind())
			*/

		}

	}

	fmt.Println("--------")
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	query := string(body)
	result := executeQuery(query, schema)
	//customProcess(schema, query)
	json.NewEncoder(w).Encode(result)
}

func slog(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txn := app.StartTransaction("query", nil, nil)
		defer txn.End()
		h.ServeHTTP(w, r) // call original
	})
}

var schema graphql.Schema
var data map[string]Car
var app newrelic.Application

func importData() {
	data = make(map[string]Car)
	data["b500"] = Car{Name: "b500", Company: "bmw", Colour: "blue", Abs: true, Price: 3000000}
	data["indigo"] = Car{Name: "indigo", Company: "tata", Colour: "black", Abs: false, Price: 1000000}
	data["swift"] = Car{Name: "swift", Company: "maruti", Colour: "grey", Abs: true, Price: 500000}
}
func initSchema() graphql.Schema {
	router := NewRouter()
	gen := tools.NewGenerator(router)
	query := gen.GenerateObject(Query{})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: query,
	})
	if err != nil {
		log.Fatal(err)
	}
	return schema
}
func initNewRelic() newrelic.Application {
	const licKey = "newrelickey"
	config := newrelic.NewConfig("tut3", licKey)
	app, err := newrelic.NewApplication(config)
	if err != nil {
		panic(err)
	}
	return app
}

func dummyQuery() {
	q := `query {
				user{
					name
				}
			}`
	res := executeQuery(q, schema)
	json.NewEncoder(os.Stdout).Encode(res)
}
func init() {
	schema = initSchema()
	app = initNewRelic()
	importData()
	//	dummyQuery()
}

func main() {
	http.HandleFunc("/graphql", slog(queryHandler))
	http.ListenAndServe(":8080", nil)
}
