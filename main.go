package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Pretty     bool
	GraphiQL   bool
	Playground bool

	Host string
	Port int

	Root string

	Projects string
}

func main() {
	configPath := pflag.StringP("config", "c", "config.yaml", "path to config file")
	pflag.Parse()

	// read yaml config
	config := Config{}
	file, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("cannot read yaml config: %v", err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("cannot parse yaml config: %v", err)
	}

	// create graphql schema
	projects := Projects{path: config.Projects}
	fprojects := projects.Field()
	fields := graphql.Fields{
		"projects": &fprojects,
	}
	root := graphql.ObjectConfig{Name: config.Root, Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(root)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// register handler
	handler := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     config.Pretty,
		GraphiQL:   config.GraphiQL,
		Playground: config.Playground,
	})
	http.Handle("/", handler)

	// run server
	fmt.Println("listening")
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	err = http.ListenAndServe(addr, nil)
	log.Fatalf("server has stopped: %v", err)
}
