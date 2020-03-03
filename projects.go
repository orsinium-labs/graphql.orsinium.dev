package main

import (
	"github.com/graphql-go/graphql"
)

func ProjectHandler(params graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}

func ProjectField() graphql.Field {
	projectType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Product",
			Fields: graphql.Fields{
				"language": &graphql.Field{
					Type: graphql.String,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"info": &graphql.Field{
					Type: graphql.String,
				},
				"link": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	return graphql.Field{
		Type:        graphql.NewList(projectType),
		Description: "Get open-source projects list",
		Resolve:     ProjectHandler,
	}
}
