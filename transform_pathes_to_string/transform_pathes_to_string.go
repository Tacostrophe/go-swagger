package transform_pathes_to_string

import (
	"fmt"
	"strings"

	S "github.com/Tacostrophe/go-swagger/structs"
)

func TransformPathesToString(pathesMethodes []S.PathMethod) (string, error) {
	var pathesStrBuilder strings.Builder
	for idx, pathMethod := range pathesMethodes {
		currentPath := pathMethod.Path

		// don't write path if previous is the same
		if isFirstElement := idx == 0; !isFirstElement && pathesMethodes[idx].Path == pathesMethodes[idx-1].Path {
			currentPath = strings.Repeat(" ", len(currentPath))
		}

		currentPathStr := fmt.Sprintf(
			"%3d. %s %s",
			idx,
			currentPath,
			strings.ToUpper(pathMethod.Method),
		)

		if isLastElement := idx == len(pathesMethodes)-1; !isLastElement {
			currentPathStr += "\n"
		}
		pathesStrBuilder.Write([]byte(currentPathStr))
	}

	return pathesStrBuilder.String(), nil
}
