package main

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/ian-kent/go-log/log"
	"net/http"
)

func serveGraphQL(port int) *http.Server {
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Error("failed to create new schema, error: %v", err)
	}
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	// TODO: Use the port from config
	server := &http.Server{Addr: ":8080"}
	http.Handle("/graphql", h)
	go server.ListenAndServe()
	log.Info("Serving at %s", server.Addr)
	return server
}
