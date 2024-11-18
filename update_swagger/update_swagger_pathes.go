package update_swagger

import (
	S "github.com/Tacostrophe/go-swagger/structs"
)

func filterPathes(swaggerPathes map[string]interface{}, pathesToKeep []S.PathMethod) (filteredSwaggerPathes map[string]map[string]interface{}) {
	filteredSwaggerPathes = make(map[string]map[string]interface{})
	for _, pathToKeep := range pathesToKeep {
		_, ok := filteredSwaggerPathes[pathToKeep.Path]
		if !ok {
			filteredSwaggerPathes[pathToKeep.Path] = map[string]interface{}{}
		}
		currentPathMethod := swaggerPathes[pathToKeep.Path].(map[string]interface{})
		filteredSwaggerPathes[pathToKeep.Path][pathToKeep.Method] = currentPathMethod[pathToKeep.Method]
	}

	return
}

func filterTags(swaggerTags []interface{}, swaggerPathes map[string]map[string]interface{}) (filteredSwaggerTags []map[string]interface{}) {
	tagsMap := make(map[string]map[string]interface{})
	for _, tag := range swaggerTags {
		tagName, ok := tag.(map[string]interface{})["name"]
		if ok {
			tagsMap[tagName.(string)] = tag.(map[string]interface{})
		}
	}

	tagsToKeepMap := make(map[string]map[string]interface{})
	for _, path := range swaggerPathes {
		for _, method := range path {
			methodTagsNames := method.(map[string]interface{})["tags"].([]interface{})
			for _, tagName := range methodTagsNames {
				tag, ok := tagsMap[tagName.(string)]
				if ok {
					tagsToKeepMap[tagName.(string)] = tag
				}
			}
		}
	}

	filteredSwaggerTags = make([]map[string]interface{}, 0, len(tagsToKeepMap))
	for _, tag := range tagsToKeepMap {
		filteredSwaggerTags = append(filteredSwaggerTags, tag)
	}

	return
}

func UpdateSwagger(swagger map[string]interface{}, pathesToKeep []S.PathMethod) (map[string]interface{}, error) {
	incomeSwaggerPathes := swagger["paths"].(map[string]interface{})
	swaggerPathes := filterPathes(incomeSwaggerPathes, pathesToKeep)
	swagger["paths"] = swaggerPathes

	incomeSwaggerTags := swagger["tags"].([]interface{})
	if len(incomeSwaggerTags) > 0 {
		swaggerTags := filterTags(incomeSwaggerTags, swaggerPathes)
		swagger["tags"] = swaggerTags
	}

	return swagger, nil
}
