package extract_pathes

import (
	S "do-swagger/structs"
	"reflect"
	"testing"
)

func TestExtractPathes(t *testing.T) {
	postOnlyMethod := map[string]map[string]interface{}{
		"post": {
			"parameters": []interface{}{},
			"responses": map[string]map[string]string{
				"200": {
					"description": "OK",
				},
			},
		},
	}

	oneEndpointMap := map[string]interface{}{
		"path/1": postOnlyMethod,
	}

	cases := []struct {
		in   map[string]interface{}
		want []S.PathMethod
	}{
		{
			map[string]interface{}{
				"paths": oneEndpointMap,
			},
			[]S.PathMethod{
				{
					Path:   "path/1",
					Method: "post",
				},
			},
		},
	}

	for _, currentCase := range cases {
		got, _ := ExtractPathes(currentCase.in)

		if !reflect.DeepEqual(got, currentCase.want) {
			t.Errorf("Testing ExtractPathes(%+v)\ngot %+v\nwant %+v", currentCase.in, got, currentCase.want)
		} else {
			t.Log("Passed")
		}
	}
}
