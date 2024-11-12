package main

import (
	"log"
	"os"

	EP "github.com/Tacostrophe/go-swagger/extract_pathes"
	IC "github.com/Tacostrophe/go-swagger/init_context"
	RS "github.com/Tacostrophe/go-swagger/read_swagger"
)

func main() {
	// get somehow path/to/swagger.json
	// get somehow name of a result file
	ctx, err := IC.InitContext(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// read swagger
	swagger, err := RS.ReadSwagger(ctx.IncomeSwaggerPath)
	if err != nil {
		log.Fatal(err)
	}

	// extract pathes from swagger
	pathesMethodes, err := EP.ExtractPathes(swagger)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("pathesMethodes: %+v", pathesMethodes)

	// ask which path to keep

	// filter pathes in swagger

	// write swagger into a file
}
