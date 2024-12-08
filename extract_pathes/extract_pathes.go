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
		for methodName, method := range methods {
			currentPathMethod := S.PathMethod{
				Path:   path,
				Method: methodName,
			}
			if tags, hasTags := method.(map[string]interface{})["tags"]; hasTags {
				if tagsArr := tags.([]interface{}); len(tagsArr) > 0 {
					firstTag := tagsArr[0].(string)
					currentPathMethod.FirstTag = firstTag
				}
			}
			pathesMethods = append(pathesMethods, currentPathMethod)
		}
	}

	sort.Slice(pathesMethods, func(i, j int) bool {
		if pathesMethods[i].FirstTag == pathesMethods[j].FirstTag {
			if pathesMethods[i].Path == pathesMethods[j].Path {
				return pathesMethods[i].Method < pathesMethods[j].Method
			}
			return pathesMethods[i].Path < pathesMethods[j].Path
		}
		return pathesMethods[i].FirstTag < pathesMethods[j].FirstTag
	})

	return pathesMethods, nil
}
