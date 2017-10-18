#!/bin/sh
curl 'http://127.0.0.1:8080/graphql' -d '
 query {
 	user{
		name
	}
	car(name:"swift"){
		name
		price
		colour
	}
}'

	       
