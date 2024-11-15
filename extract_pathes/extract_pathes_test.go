package extract_pathes

import (
	"reflect"
	"testing"

	S "github.com/Tacostrophe/go-swagger/structs"
)

func TestExtractPathes(t *testing.T) {
	oneEndpointMap := map[string]map[string]interface{}{
		"path/1": {
			"post": map[string]interface{}{
				"parameters": []interface{}{},
				"responses": map[string]map[string]string{
					"200": {
						"description": "OK",
					},
				},
			},
		},
	}

	coupleEndpointMap := map[string]map[string]interface{}{
		"path/2": {
			"post": map[string]interface{}{
				"parameters": []interface{}{},
				"responses": map[string]map[string]string{
					"200": {
						"description": "OK",
					},
				},
			},
		},
		"path/1": {
			"post": map[string]interface{}{
				"parameters": []interface{}{},
				"responses": map[string]map[string]string{
					"200": {
						"description": "OK",
					},
				},
			},
			"get": map[string]interface{}{
				"parameters": []interface{}{},
				"responses": map[string]map[string]string{
					"200": {
						"description": "OK",
					},
				},
			},
		},
	}

	cases := []struct {
		in   S.Swagger
		want []S.PathMethod
	}{
		{
			S.Swagger{
				Paths: map[string]map[string]interface{}{},
			},
			[]S.PathMethod{},
		},
		{
			S.Swagger{
				Paths: oneEndpointMap,
			},
			[]S.PathMethod{
				{
					Path:   "path/1",
					Method: "post",
				},
			},
		},
		{
			S.Swagger{
				Paths: coupleEndpointMap,
			},
			[]S.PathMethod{
				{
					Path:   "path/1",
					Method: "get",
				},
				{
					Path:   "path/1",
					Method: "post",
				},
				{
					Path:   "path/2",
					Method: "post",
				},
			},
		},
	}

	for _, currentCase := range cases {
		got, _ := ExtractPathes(currentCase.in)

		wantLen := len(currentCase.want)
		gotLen := len(got)

		if gotLen != wantLen {
			t.Errorf("Expected array length isn't equal result\ngot %d\nwant %d", gotLen, wantLen)
		}

		if gotLen != 0 && !reflect.DeepEqual(got, currentCase.want) {
			t.Errorf("Testing ExtractPathes(%+v)\ngot %+v\nwant %+v", currentCase.in, got, currentCase.want)
		} else {
			t.Log("Passed")
		}
	}
}
