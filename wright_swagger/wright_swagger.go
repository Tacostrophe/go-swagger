package wright_swagger

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

func WrightSwagger(swagger map[string]interface{}) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Input file/path/to/save.json: ")

	rawInput, _ := reader.ReadString('\n')
	resultFileName := strings.Replace(rawInput, "\n", "", -1)

	if resultFileName == "" {
		currentTime := time.Now()
		resultFileName = fmt.Sprintf("swagger%s.json", currentTime.Format(time.RFC3339))
	}

	if !strings.HasSuffix(resultFileName, ".json") {
		return "", errors.New("result file must have json extension")
	}

	if strings.Contains(resultFileName, "~") {
		homeDirPath, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		resultFileName = strings.Replace(resultFileName, "~", homeDirPath, 1)
	}

	swaggerString, err := json.Marshal(swagger)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(resultFileName, swaggerString, 0644)
	if err != nil {
		return "", err
	}

	return resultFileName, nil
}
