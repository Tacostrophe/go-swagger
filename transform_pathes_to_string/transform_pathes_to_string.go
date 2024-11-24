package transform_pathes_to_string

import (
	"errors"
	"fmt"
	"strings"

	S "github.com/Tacostrophe/go-swagger/structs"
)

func TransformPathesToString(pathesMethodes []S.PathMethod) (string, error) {
	pathesAmount := len(pathesMethodes)
	if pathesAmount == 0 {
		return "", errors.New("swagger must have at least one path")
	}
	rowsAmount := 50
	maxRowLengths := make([]int, pathesAmount/rowsAmount)
	minIndentation := 5
	if pathesAmount < rowsAmount {
		rowsAmount = pathesAmount
	}

	pathesRows := make([]string, rowsAmount)

	for idx, pathMethod := range pathesMethodes {
		currentPath := pathMethod.Path

		// don't write path if previous is the same
		if isFirstElement := idx == 0; !isFirstElement && currentPath == pathesMethodes[idx-1].Path {
			currentPath = strings.Repeat(" ", len(currentPath))
		}

		currentPathStr := fmt.Sprintf(
			"%3d. %s %s",
			idx,
			currentPath,
			strings.ToUpper(pathMethod.Method),
		)

		rowIdx := idx % rowsAmount
		columnIdx := idx / rowsAmount
		if columnIdx > 0 && columnIdx <= len(maxRowLengths) {
			indentation := strings.Repeat(" ", maxRowLengths[columnIdx-1]-len(pathesRows[rowIdx]))
			currentPathStr = indentation + currentPathStr
		}
		pathesRows[rowIdx] += currentPathStr
		currentRowLen := len(pathesRows[rowIdx]) + minIndentation
		if columnIdx < len(maxRowLengths) && currentRowLen > maxRowLengths[columnIdx] {
			maxRowLengths[columnIdx] = currentRowLen
		}
	}

	return strings.Join(pathesRows, "\n"), nil
}
