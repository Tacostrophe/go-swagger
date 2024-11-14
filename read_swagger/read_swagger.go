package read_swagger

import (
	"encoding/json"
	"io/ioutil"
	"os"

	S "github.com/Tacostrophe/go-swagger/structs"
)

func ReadSwagger(path string) (S.Swagger, map[string]interface{}, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return S.Swagger{}, map[string]interface{}{}, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return S.Swagger{}, map[string]interface{}{}, err
	}

	var swaggerPathes S.Swagger
	err = json.Unmarshal([]byte(byteValue), &swaggerPathes)
	if err != nil {
		return S.Swagger{}, map[string]interface{}{}, err
	}

	var swagger map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &swagger)
	if err != nil {
		return S.Swagger{}, map[string]interface{}{}, err
	}
	return swaggerPathes, swagger, nil
}
