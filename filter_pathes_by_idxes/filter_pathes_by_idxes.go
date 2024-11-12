package filter_pathes_by_idxes

import (
	"errors"
	"fmt"

	S "github.com/Tacostrophe/go-swagger/structs"
)

func FilterPathesByIdxes(pathes []S.PathMethod, pathesToKeep []int) ([]S.PathMethod, error) {
	var filteredPathes []S.PathMethod
	for _, idx := range pathesToKeep {
		if idx >= len(pathes) {
			errorMessage := fmt.Sprintf("Provided index is out of range: %d", idx)
			return []S.PathMethod{}, errors.New(errorMessage)
		}
		filteredPathes = append(filteredPathes, pathes[idx])
	}

	return filteredPathes, nil
}
