package router

import "github.com/prasannamahajan/go-graphql-tools"
import "github.com/prasannamahajan/graphql-newrelic/schema"
import "github.com/prasannamahajan/graphql-newrelic/data"

func NewRouter() *tools.Router {
	router := tools.NewRouter()
	router.Query("Query.User", func() (schema.User, error) {
		return data.GetUser(), nil
	})
	router.Query("Query.Car", func(q schema.Query, args schema.QueryCarArgs) (interface{}, error) {
		return data.GetCar(args.Name), nil
	})
	router.Query("Car.Price", func(c schema.Car, args schema.CarPriceArgs) (interface{}, error) {
		if args.Unit != nil {
			if *args.Unit == "USD" {
				return c.Price / 70, nil
			}
		}
		return c.Price, nil
	})
	return router
}
