package update_swagger

import (
	"errors"
	"sort"

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

	sort.Slice(filteredSwaggerTags, func(i, j int) bool {
		return filteredSwaggerTags[i]["name"].(string) < filteredSwaggerTags[j]["name"].(string)
	})

	return
}

func getRefsFromMap(currentMap map[string]interface{}) (refs map[string]bool) {
	refs = make(map[string]bool)
	if len(currentMap) == 0 {
		return
	}

	for key, value := range currentMap {
		if key == "$ref" {
			refs[value.(string)] = true
			continue
		}

		subMap, isMap := value.(map[string]interface{})
		if !isMap {
			continue
		}

		subMapRefs := getRefsFromMap(subMap)
		for refName := range subMapRefs {
			refs[refName] = true
		}
	}

	return
}

func getRefsFromPathes(pathes map[string]map[string]interface{}) (refs map[string]bool) {
	refs = make(map[string]bool)
	for _, methods := range pathes {
		pathRefs := getRefsFromMap(methods)
		for refName := range pathRefs {
			refs[refName] = true
		}
	}

	return
}

func filterComponentsSchemas(swaggerComponentsSchemas map[string]interface{}, swaggerPathes map[string]map[string]interface{}) (filteredSwaggerComponentsSchemas map[string]interface{}) {
	filteredSwaggerComponentsSchemas = make(map[string]interface{})
	refs := getRefsFromPathes(swaggerPathes)
	for schemaName, schema := range swaggerComponentsSchemas {
		refName := "#/components/schemas/" + schemaName
		if _, hasRef := refs[refName]; hasRef {
			filteredSwaggerComponentsSchemas[schemaName] = schema
		}
	}
	return
}

func UpdateSwagger(swagger map[string]interface{}, pathesToKeep []S.PathMethod) (map[string]interface{}, error) {
	incomeSwaggerPathesI, ok := swagger["paths"]
	if !ok {
		return map[string]interface{}{}, errors.New("swagger to update must have pathes")
	}
	incomeSwaggerPathes, ok := incomeSwaggerPathesI.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, errors.New("didn't manage to cast swagger pathes to map")
	}
	swaggerPathes := filterPathes(incomeSwaggerPathes, pathesToKeep)
	swagger["paths"] = swaggerPathes

	incomeSwaggerTags := swagger["tags"].([]interface{})
	if len(incomeSwaggerTags) > 0 {
		swaggerTags := filterTags(incomeSwaggerTags, swaggerPathes)
		swagger["tags"] = swaggerTags
	}

	if incomeSwaggerComponents, hasComponents := swagger["components"].(map[string]interface{}); hasComponents {
		if incomeSwaggerComponentsSchemas, hasSchemas := incomeSwaggerComponents["schemas"].(map[string]interface{}); hasSchemas {
			if len(incomeSwaggerComponentsSchemas) > 0 {
				filteredSwaggerComponentsSchemas := filterComponentsSchemas(incomeSwaggerComponentsSchemas, swaggerPathes)
				incomeSwaggerComponents["schemas"] = filteredSwaggerComponentsSchemas
				swagger["components"] = incomeSwaggerComponents
			}
		}
	}

	return swagger, nil
}
