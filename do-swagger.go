package main

import (
	"log"

	EP "github.com/Tacostrophe/go-swagger/extract_pathes"
	RS "github.com/Tacostrophe/go-swagger/read_swagger"
)

func main() {
	// get somehow path/to/swagger.json
	// get somehow name of a result file
	path := "/home/tacostrophe/Downloads/swagger.json"
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
