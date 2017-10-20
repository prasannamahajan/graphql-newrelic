package main

import (
	"encoding/json"
	"fmt"
	"github.com/prasannamahajan/go-graphql-tools"
	"github.com/prasannamahajan/graphql"
	"github.com/prasannamahajan/graphql-newrelic/data"
	"github.com/prasannamahajan/graphql-newrelic/relicconf"
	"github.com/prasannamahajan/graphql-newrelic/router"
	"github.com/prasannamahajan/graphql-newrelic/schema"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func Usage() {
	fmt.Printf("\nUsage : %s <appname> <lickey>\n", os.Args[0])
	os.Exit(1)
}
func parseArgs() (string, string) {
	if len(os.Args) != 3 {
		Usage()
	}
	return os.Args[1], os.Args[2]
}

func init() {
	Schema = InitSchema()
	appname, licKey := parseArgs()
	relicconf.InitNewRelic(appname, licKey)
	data.ImportData()
}

func main() {
	http.HandleFunc("/graphql", queryHandler)
	log.Println("Starting Server")
	http.ListenAndServe(":8080", nil)
}
