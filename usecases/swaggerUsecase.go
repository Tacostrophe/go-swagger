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
