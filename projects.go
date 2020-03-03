package main

import (
	"fmt"
	"io/ioutil"

	"github.com/graphql-go/graphql"
	"gopkg.in/yaml.v2"
)

type Project struct {
	Name     string
	Info     string
	Link     string
	Language string
}

type Language struct {
	Name     string
	Projects []Project `yaml:"items"`
}

type Projects struct {
	path  string
	cache []Project
}

func (pr Projects) read() ([]Project, error) {
	langs := []Language{}
	file, err := ioutil.ReadFile(pr.path)
	if err != nil {
		return nil, fmt.Errorf("cannot read yaml file for projects: %w", err)
	}
	err = yaml.Unmarshal(file, &langs)
	if err != nil {
		return nil, fmt.Errorf("cannot parse yaml file for projects: %w", err)
	}

	result := make([]Project, 0)
	for _, lang := range langs {
		for _, project := range lang.Projects {
			project.Language = lang.Name
			result = append(result, project)
		}
	}

	return result, nil
}

func (pr *Projects) Handle(params graphql.ResolveParams) (interface{}, error) {
	var err error
	if pr.cache == nil {
		pr.cache, err = pr.read()
	}
	return pr.cache, err
}

func (pr *Projects) Field() graphql.Field {
	projectType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Project",
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
		Resolve:     pr.Handle,
	}
}
