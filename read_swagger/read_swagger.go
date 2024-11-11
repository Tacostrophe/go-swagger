package read_swagger

import (
	"encoding/json"
	"io/ioutil"
	"os"

	S "github.com/Tacostrophe/go-swagger/structs"
)

func ReadSwagger(path string) (S.Swagger, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return S.Swagger{}, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return S.Swagger{}, err
	}

	var swagger S.Swagger
	err = json.Unmarshal([]byte(byteValue), &swagger)
	if err != nil {
		return S.Swagger{}, err
	}
	return swagger, nil
}
