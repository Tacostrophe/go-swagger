package extract_pathes

import (
	"errors"
	"sort"

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

	sort.Slice(pathesMethods, func(i, j int) bool {
		return pathesMethods[i].Path < pathesMethods[j].Path || pathesMethods[i].Method < pathesMethods[j].Method
	})

	return pathesMethods, nil
}
