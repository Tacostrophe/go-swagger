package extract_pathes

import (
	S "do-swagger/structs"
	"reflect"
	"sort"
	"testing"
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
		"path/1": {
			"get": map[string]interface{}{
				"parameters": []interface{}{},
				"responses": map[string]map[string]string{
					"200": {
						"description": "OK",
					},
				},
			},
			"post": map[string]interface{}{
				"parameters": []interface{}{},
				"responses": map[string]map[string]string{
					"200": {
						"description": "OK",
					},
				},
			},
		},
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
	}

	cases := []struct {
		in   map[string]map[string]map[string]interface{}
		want []S.PathMethod
	}{
		{
			map[string]map[string]map[string]interface{}{
				"paths": {},
			},
			[]S.PathMethod{},
		},
		{
			map[string]map[string]map[string]interface{}{
				"paths": oneEndpointMap,
			},
			[]S.PathMethod{
				{
					Path:   "path/1",
					Method: "post",
				},
			},
		},
		{
			map[string]map[string]map[string]interface{}{
				"paths": coupleEndpointMap,
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

		comparePathMethods := func(i, j int) bool {
			return got[i].Path < got[j].Path && got[i].Method < got[j].Method
		}
		sort.Slice(got[:], comparePathMethods)
		sort.Slice(currentCase.want[:], comparePathMethods)

		if gotLen != 0 && !reflect.DeepEqual(got, currentCase.want) {
			t.Errorf("Testing ExtractPathes(%+v)\ngot %+v\nwant %+v", currentCase.in, got, currentCase.want)
		} else {
			t.Log("Passed")
		}
	}
}
