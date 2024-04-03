package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"mywon/students_reports/graph"
	"mywon/students_reports/repository"
	"mywon/students_reports/utility"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func init() {
	environment := ""
	if len(os.Args) < 2 {
		log.Println("missing application environment method. Usage: go run server.go {dev | uat | prod} ")
		os.Exit(1)
	}
	environment = os.Args[1]

	fmt.Println("#####",environment)
	if environment == "docker"{
		environment = "dev"
	}
	if environment == "dev" {
		err := godotenv.Load(".dev-env")
		if err != nil {
			log.Println("failed to load .dev-env file")
		}
	} else {
		log.Println("Environment '" + environment + "' not available. Usage: go run server.go {dev} ")
		os.Exit(1)
	}
	repository.InitDbPool()
}

func main() {
	defer utility.RecoverFromPanic("server.go")
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("Invalid PORT, application exiting...")
		os.Exit(1)
	}
	log.Printf("Port is %v\n", port)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
