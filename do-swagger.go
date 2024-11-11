package main

import (
	"log"
	"os"
	"strings"

	EP "github.com/Tacostrophe/go-swagger/extract_pathes"
	RS "github.com/Tacostrophe/go-swagger/read_swagger"
)

func main() {
	// get somehow path/to/swagger.json
	// get somehow name of a result file
	if len(os.Args) < 2 {
		log.Fatal("path/to/json as first argument is required")
	}
	path := os.Args[1]
	if !strings.HasSuffix(path, ".json") {
		log.Fatal("file must be a json")
	}

	swagger, err := RS.ReadSwagger(path)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("swagger:", swagger)

	// extract pathes
	pathesMethodes, err := EP.ExtractPathes(swagger)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("pathesMethodes:", pathesMethodes)
	// ask which path to keep

	// filter pathes in swagger

	// write swagger into a file
}
