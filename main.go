package main

import (
	"encoding/json"
	"github.com/arvitaly/go-graphql-tools"
	"github.com/arvitaly/graphql"
	"graphql-newrelic/data"
	"graphql-newrelic/relicconf"
	"graphql-newrelic/router"
	"graphql-newrelic/schema"
	"io/ioutil"
	"log"
	"net/http"
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
	log.Println("query is : ", query)
	result := graphql.Do(graphql.Params{
		RequestString: query,
		Schema:        schema,
	},
	)
	if len(result.Errors) > 0 {
		log.Println("some error occured")
	}
	return result
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	query := string(body)
	result := executeQuery(query, Schema)
	json.NewEncoder(w).Encode(result)
}

var Schema graphql.Schema

func init() {
	Schema = InitSchema()
	relicconf.InitNewRelic()
	data.ImportData()
}

func main() {
	http.HandleFunc("/graphql", queryHandler)
	log.Println("Starting Server")
	http.ListenAndServe(":8080", nil)
}
