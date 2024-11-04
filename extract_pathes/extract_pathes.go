package extract_pathes

import (
	S "do-swagger/structs"
	"errors"
	"fmt"
)

func ExtractPathes(swagger map[string]interface{}) ([]S.PathMethod, error) {
	pathes, ok := swagger["paths"]
	if !ok {
		return []S.PathMethod{}, errors.New("swagger must have pathes")
	}

	fmt.Printf("pathes: %v\n", pathes)
	fmt.Printf("pathes len: %d\n", len(pathes.(map[string]interface{})))

	return []S.PathMethod{}, nil
}
