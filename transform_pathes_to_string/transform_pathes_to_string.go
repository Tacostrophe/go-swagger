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
		pathesStrBuilder.Write([]byte(idxString))
		pathesStrBuilder.Write([]byte(" "))
		pathesStrBuilder.Write([]byte(strings.ToUpper(pathMethod.Method)))
		pathesStrBuilder.Write([]byte(" "))
		pathesStrBuilder.Write([]byte(pathMethod.Path))
		pathesStrBuilder.Write([]byte(" "))
		pathesStrBuilder.Write([]byte("\n"))
	}

	return pathesStrBuilder.String(), nil
}
