package main

import (
	"fmt"
	"gq_enrich/cfg"
	"gq_enrich/graph"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	cMainCfg := &cfg.Cfg{}
	err := cfg.LoadConfig(cMainCfg)
	CheckError(err)

	if cMainCfg.Port == "" {
		cMainCfg.Port = defaultPort
	}

	dBase := graph.Connect(cMainCfg.PostgresURL)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: dBase}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cMainCfg.Port)
	log.Fatal(http.ListenAndServe(":"+cMainCfg.Port, nil))
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
