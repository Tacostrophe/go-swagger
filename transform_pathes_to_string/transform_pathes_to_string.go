package transform_pathes_to_string

import (
	"strconv"
	"strings"

	S "github.com/Tacostrophe/go-swagger/structs"
)

func TransformPathesToString(pathesMethodes []S.PathMethod) (string, error) {
	var pathesStrBuilder strings.Builder
	for idx, pathMethod := range pathesMethodes {
		idxString := strconv.Itoa(idx)
		currentPathStr := idxString + " " + strings.ToUpper(pathMethod.Method) + " " + pathMethod.Path + "\n"
		pathesStrBuilder.Write([]byte(currentPathStr))
	}

	return pathesStrBuilder.String(), nil
}
