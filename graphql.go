package main

import (
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/ian-kent/go-log/log"
	"net/http"
	"time"
)

type ReleaseResponse struct {
	Successful bool
	Message    string
	Id         int
}

var releaseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Release",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"version": &graphql.Field{
				Type: graphql.String,
			},
			"dateSubmitted": &graphql.Field{
				Type: graphql.DateTime,
			},
			"released": &graphql.Field{
				Type: graphql.Boolean,
			},
			"dateReleased": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

var releaseResponseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ReleaseReponse",
		Fields: graphql.Fields{
			"successful": &graphql.Field{
				Type: graphql.Boolean,
			},
			"id": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"Release": &graphql.Field{
				Type: releaseType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					idQuery, ok := p.Args["id"].(int)
					if !ok {
						return nil, nil
					}
					release, ok := Releases[idQuery]
					if !ok {
						return nil, errors.New(fmt.Sprintf("Failed to find release with ID %d", idQuery))
					}
					return release, nil
				},
			},
			"Releases": &graphql.Field{
				Type: graphql.NewList(releaseType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					releases := []Release{}
					for _, value := range Releases {
						releases = append(releases, value)
					}
					return releases, nil
				},
			},
		},
	},
)

var mutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"newRelease": &graphql.Field{
				Type:        releaseResponseType,
				Description: "Create a release request",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"version": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					name, _ := params.Args["name"].(string)
					version, _ := params.Args["version"].(string)

					newId := len(Releases)
					release := Release{
						ID:            newId,
						Name:          name,
						Version:       version,
						DateSubmitted: time.Now(),
					}
					err := AddRelease(release)
					if err != nil {
						response := ReleaseResponse{
							Successful: false,
							Message:    fmt.Sprintf("%s", err),
						}
						return response, nil
					}
					response := ReleaseResponse{
						Successful: true,
						Id:         newId,
					}
					return response, nil
				},
			},
		},
	},
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

func serveGraphQL(port int) *http.Server {
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
