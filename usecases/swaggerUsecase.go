package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type (
	swaggerFromFileV1 struct {
		filePath string
		swagger  swagger

		pathes []PathMethod
	}

	swagger map[string]interface{}

	PathMethod struct {
		Path     string
		Method   string
		FirstTag string
	}
)

func NewSwaggerFromFileV1() (SwaggerUsecase, error) {
	return &swaggerFromFileV1{}, nil
}

func (u *swaggerFromFileV1) Init(filePath string) error {
	swagger, err := readSwagger(filePath)
	if err != nil {
		return err
	}

	pathes, err := extractPathes(swagger)
	if err != nil {
		return err
	}

	u.swagger = swagger
	u.pathes = pathes

	return nil
}

func readSwagger(filePath string) (swagger, error) {
	if strings.HasPrefix(filePath, "~") {
		homeDirPath, err := os.UserHomeDir()
		if err != nil {
			return swagger{}, err
		}

		filePath = strings.Replace(filePath, "~", homeDirPath, 1)
	}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return swagger{}, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return swagger{}, err
	}

	var swaggerContent swagger
	err = json.Unmarshal([]byte(byteValue), &swaggerContent)
	if err != nil {
		return swagger{}, err
	}
	return swaggerContent, nil
}

func extractPathes(swagger swagger) ([]PathMethod, error) {
	pathesVal, exists := swagger["paths"]
	if !exists {
		errMessage := fmt.Sprintf("swagger has no paths in it: %t", exists)
		return []PathMethod{}, errors.New(errMessage)
	}

	pathes := pathesVal.(map[string]interface{})
	// if !ok {
	// 	return []pathMethod{}, errors.New("incorrect structure of \"pathes\"")
	// }
	if len(pathes) == 0 {
		return []PathMethod{}, errors.New("swagger has no paths in it")
	}

	var pathesMethods []PathMethod
	for path, methods := range pathes {
		for methodName, method := range methods.(map[string]interface{}) {
			currentPathMethod := PathMethod{
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

func (u *swaggerFromFileV1) GetFilteredPathes(filter string) []PathMethod {
	if filter == "" {
		return u.pathes
	}

	filteredPathes := make([]PathMethod, 0, len(u.pathes))
	for _, path := range u.pathes {
		if strings.Contains(path.Path, filter) || strings.Contains(path.Method, filter) || strings.Contains(path.FirstTag, filter) {
			filteredPathes = append(filteredPathes, path)
		}
	}

	return filteredPathes
}

func (u *swaggerFromFileV1) filterComponentsSchemas(swaggerComponentsSchemas map[string]interface{}, swaggerPathes map[string]map[string]interface{}) (filteredSwaggerComponentsSchemas map[string]interface{}) {
	return
}

func (u *swaggerFromFileV1) filterTags(swaggerPathes map[string]map[string]interface{}) (filteredSwaggerTags []map[string]interface{}) {
	swaggerTags := u.swagger["tags"].([]interface{})
	if len(swaggerTags) == 0 {
		return
	}

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
			methodTagsNames, ok := method.(map[string]interface{})["tags"].([]interface{})
			if !ok {
				continue
			}
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

func (u *swaggerFromFileV1) filterPathes(pathesToKeep []PathMethod) (map[string]map[string]interface{}, error) {
	incomeSwaggerPathesI, ok := u.swagger["paths"]
	if !ok {
		return nil, errors.New("swagger to update must have pathes")
	}
	swaggerPathes, ok := incomeSwaggerPathesI.(map[string]interface{})
	if !ok {
		return nil, errors.New("didn't manage to cast swagger pathes to map")
	}

	filteredSwaggerPathes := make(map[string]map[string]interface{})
	for _, pathToKeep := range pathesToKeep {
		_, ok := filteredSwaggerPathes[pathToKeep.Path]
		if !ok {
			filteredSwaggerPathes[pathToKeep.Path] = map[string]interface{}{}
		}
		currentPathMethod := swaggerPathes[pathToKeep.Path].(map[string]interface{})
		filteredSwaggerPathes[pathToKeep.Path][pathToKeep.Method] = currentPathMethod[pathToKeep.Method]
	}

	return filteredSwaggerPathes, nil
}

func (u *swaggerFromFileV1) UpdateSwagger(pathesToKeep []PathMethod) error {
	swagger := u.swagger
	swaggerPathes, err := u.filterPathes(pathesToKeep)
	if err != nil {
		return err
	}

	swaggerTags := u.filterTags(swaggerPathes)

	if incomeSwaggerComponents, hasComponents := swagger["components"].(map[string]interface{}); hasComponents {
		if incomeSwaggerComponentsSchemas, hasSchemas := incomeSwaggerComponents["schemas"].(map[string]interface{}); hasSchemas {
			if len(incomeSwaggerComponentsSchemas) > 0 {
				filteredSwaggerComponentsSchemas := u.filterComponentsSchemas(incomeSwaggerComponentsSchemas, swaggerPathes)
				incomeSwaggerComponents["schemas"] = filteredSwaggerComponentsSchemas
				swagger["components"] = incomeSwaggerComponents
			}
		}
	}

	swagger["paths"] = swaggerPathes
	if swaggerTags != nil {
		swagger["tags"] = swaggerTags
	}
	u.swagger = swagger

	return nil
}
