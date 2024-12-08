package transform_pathes_to_string

import (
	"errors"
	"fmt"
	"strings"

	S "github.com/Tacostrophe/go-swagger/structs"
)

func groupPathesByFirstTag(pathesMethodes []S.PathMethod) [][]S.PathMethod {
	// expected income array to be sorted by first tag
	groupedPathesMethodes := [][]S.PathMethod{}
	if len(pathesMethodes) == 0 {
		return groupedPathesMethodes
	}

	currentFirstTag := pathesMethodes[0].FirstTag
	currentFirstTagStartIdx := 0

	for idx := 1; idx < len(pathesMethodes); idx++ {
		if pathMethod := pathesMethodes[idx]; currentFirstTag != pathMethod.FirstTag {
			tagName := currentFirstTag
			if tagName == "" {
				tagName = "*Tagless*"
			}
			groupedPathesMethodes = append(groupedPathesMethodes, pathesMethodes[currentFirstTagStartIdx:idx])
			currentFirstTag = pathMethod.FirstTag
			currentFirstTagStartIdx = idx
		}
	}
	groupedPathesMethodes = append(groupedPathesMethodes, pathesMethodes[currentFirstTagStartIdx:])
	return groupedPathesMethodes
}

func TransformPathesToString(pathesMethodes []S.PathMethod) (string, error) {
	pathesAmount := len(pathesMethodes)
	if pathesAmount == 0 {
		return "", errors.New("swagger must have at least one path")
	}
	columnsAmount := 3
	minIndentation := 5

	groupedPathes := groupPathesByFirstTag(pathesMethodes)
	rowsAmount := 0
	staticCharsAmountInPathString := 6
	maxRowLengths := make([]int, columnsAmount-1)
	for _, tagsPathes := range groupedPathes {
		currentTagMaxRowsPerColumn := (len(tagsPathes) / columnsAmount)
		if len(tagsPathes)%columnsAmount > 0 {
			currentTagMaxRowsPerColumn += 1
		}
		rowsAmount += 1 + currentTagMaxRowsPerColumn
		for tagPathIdx, tagPath := range tagsPathes {
			columnIdx := tagPathIdx / currentTagMaxRowsPerColumn
			currentRowLen := len(tagPath.Method) + len(tagPath.Path) + staticCharsAmountInPathString + minIndentation
			if columnIdx < len(maxRowLengths) && currentRowLen > maxRowLengths[columnIdx] {
				maxRowLengths[columnIdx] = currentRowLen
			}
		}
	}
	for mrlIdx := range maxRowLengths {
		if mrlIdx == 0 {
			continue
		}

		maxRowLengths[mrlIdx] += maxRowLengths[mrlIdx-1]
	}

	pathesRows := make([]string, rowsAmount)

	pathIdx := 0
	rowIdx := 0
	for _, tagPathes := range groupedPathes {
		tag := tagPathes[0].FirstTag
		if tag == "" {
			tag = "*Tagless*"
		}
		pathesRows[rowIdx] = tag
		rowIdx += 1
		rowsPerColumn := (len(tagPathes) / columnsAmount)
		if len(tagPathes)%columnsAmount > 0 {
			rowsPerColumn += 1
		}

		for tagPathIdx, pathMethod := range tagPathes {
			currentPath := pathMethod.Path

			// don't write path if previous is the same
			if isFirstElement := tagPathIdx%rowsPerColumn == 0; !isFirstElement && currentPath == tagPathes[tagPathIdx-1].Path {
				currentPath = strings.Repeat(" ", len(currentPath))
			}

			currentPathStr := fmt.Sprintf(
				"%3d. %s %s",
				pathIdx,
				currentPath,
				strings.ToUpper(pathMethod.Method),
			)
			pathIdx += 1

			currentRowIdx := (tagPathIdx % rowsPerColumn) + rowIdx
			columnIdx := tagPathIdx / rowsPerColumn
			if columnIdx > 0 && columnIdx <= len(maxRowLengths) {
				indentation := strings.Repeat(" ", maxRowLengths[columnIdx-1]-len(pathesRows[currentRowIdx]))
				currentPathStr = indentation + currentPathStr
			}
			pathesRows[currentRowIdx] += currentPathStr
		}
		rowIdx += rowsPerColumn
	}

	return strings.Join(pathesRows, "\n"), nil
}
