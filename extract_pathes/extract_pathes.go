package extract_pathes

import (
	S "do-swagger/structs"
	"errors"
)

func ExtractPathes(swagger map[string]map[string]map[string]interface{}) ([]S.PathMethod, error) {
	pathes, ok := swagger["paths"]
	if !ok {
		return []S.PathMethod{}, errors.New("swagger must have pathes")
	}

	var pathesMethods []S.PathMethod
	for path, methods := range pathes {
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
