package update_swagger_pathes

import (
	S "github.com/Tacostrophe/go-swagger/structs"
)

func UpdateSwaggerPathes(swagger map[string]interface{}, pathesToKeep []S.PathMethod) (map[string]interface{}, error) {
	incomeSwaggerPathes := swagger["paths"].(map[string]interface{})
	swaggerPathes := make(map[string]map[string]interface{})
	for _, pathToKeep := range pathesToKeep {
		_, ok := swaggerPathes[pathToKeep.Path]
		if !ok {
			swaggerPathes[pathToKeep.Path] = map[string]interface{}{}
		}
		currentPathMethod := incomeSwaggerPathes[pathToKeep.Path].(map[string]interface{})
		swaggerPathes[pathToKeep.Path][pathToKeep.Method] = currentPathMethod[pathToKeep.Method]
	}

	swagger["paths"] = swaggerPathes
	return swagger, nil
}
