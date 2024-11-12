package request_pathes_to_keep

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func RequestPathesToKeep(pathesStr string) ([]int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Presented pathes in swagger:")
	fmt.Println(pathesStr)
	fmt.Print("Input pathes indexes to keep: ")

	rawInput, _ := reader.ReadString('\n')
	sanitizedInput := strings.Replace(rawInput, "\n", "", -1)
	indexesAsStrings := strings.Split(sanitizedInput, " ")

	var indexes []int
	for _, strVal := range indexesAsStrings {
		intVal, err := strconv.Atoi(strVal)
		if err != nil {
			return []int{}, err
		}
		indexes = append(indexes, intVal)
	}

	return indexes, nil
}
