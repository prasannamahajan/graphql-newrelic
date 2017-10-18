package main

import (
	"encoding/json"
	//"fmt"
	"fmt"
	"github.com/arvitaly/go-graphql-tools"
	"github.com/arvitaly/graphql"
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

func queryHandler(w http.ResponseWriter, r *http.Request) {
	//query := r.URL.Query().Get("query")
	fmt.Println("method :", r.Method)
	body, _ := ioutil.ReadAll(r.Body)
	query := string(body)
	result := executeQuery(query, g_schema)
	json.NewEncoder(w).Encode(result)
}

func slog(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//	txn := App.StartTransaction("query", nil, nil)
		//	defer txn.End()
		h.ServeHTTP(w, r) // call original
	})
}

var g_schema graphql.Schema

func init() {
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
	g_schema = schema
	q := `query {
				user{
					name
				}
			}`
	res := executeQuery(q, g_schema)
	json.NewEncoder(os.Stdout).Encode(res)
}

func main() {
	http.HandleFunc("/graphql", slog(queryHandler))
	http.ListenAndServe(":8080", nil)
}
