package extract_pathes

import (
	"errors"

	S "github.com/Tacostrophe/go-swagger/structs"
)

func ExtractPathes(swagger S.Swagger) ([]S.PathMethod, error) {
	var pathesMethods []S.PathMethod
	if len(swagger.Paths) == 0 {
		return []S.PathMethod{}, errors.New("swagger has no paths in it")
	}
	for path, methods := range swagger.Paths {
		for method := range methods {
			currentPathMethod := S.PathMethod{
				Path:   path,
				Method: method,
			}
			pathesMethods = append(pathesMethods, currentPathMethod)
		}
	}
	return pathesMethods, nil
}
