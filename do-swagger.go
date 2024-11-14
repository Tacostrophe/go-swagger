package main

import (
	"fmt"
	"log"
	"os"

	EP "github.com/Tacostrophe/go-swagger/extract_pathes"
	FP "github.com/Tacostrophe/go-swagger/filter_pathes_by_idxes"
	IC "github.com/Tacostrophe/go-swagger/init_context"
	RS "github.com/Tacostrophe/go-swagger/read_swagger"
	RP "github.com/Tacostrophe/go-swagger/request_pathes_to_keep"
	TS "github.com/Tacostrophe/go-swagger/transform_pathes_to_string"
	US "github.com/Tacostrophe/go-swagger/update_swagger_pathes"
)

func main() {
	// get somehow path/to/swagger.json
	// get somehow name of a result file
	ctx, err := IC.InitContext(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// read swagger
	swaggerPathes, swagger, err := RS.ReadSwagger(ctx.IncomeSwaggerPath)
	if err != nil {
		log.Fatal(err)
	}

	// extract pathes from swagger
	pathesMethodes, err := EP.ExtractPathes(swaggerPathes)
	if err != nil {
		log.Fatal(err)
	}

	// transform array of pathes to human readable list with idx
	pathesStr, err := TS.TransformPathesToString(pathesMethodes)
	if err != nil {
		log.Fatal(err)
	}

	// request pathes to keep
	pathesIdxesToKeep, err := RP.RequestPathesToKeep(pathesStr)
	if err != nil {
		log.Fatal(err)
	}

	// filter pathes in swagger
	pathesToKeep, err := FP.FilterPathesByIdxes(pathesMethodes, pathesIdxesToKeep)
	if err != nil {
		log.Fatal(err)
	}

	swagger, err = US.UpdateSwaggerPathes(swagger, pathesToKeep)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("swagger:\n%+v", swagger)

	// write swagger into a file
}
