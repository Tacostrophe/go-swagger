package update_swagger

import (
	"testing"

	S "github.com/Tacostrophe/go-swagger/structs"
)

type updateSwaggerArguments struct {
	swagger      map[string]interface{}
	pathesToKeep []S.PathMethod
}

func testSwaggerFactory() map[string]interface{} {
	return map[string]interface{}{
		"paths": map[string]interface{}{
			"/path1": map[string]interface{}{
				"requestBody": map[string]interface{}{
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"type": "array",
								"items": map[string]interface{}{
									"$ref": "#/components/schemas/ref2",
								},
							},
						},
					},
				},
				"get": map[string]interface{}{
					"tags": []interface{}{
						"tag1",
					},
				},
				"post": map[string]interface{}{
					"tags": []interface{}{
						"tag1",
					},
				},
			},
			"/path5": map[string]interface{}{
				"delete": map[string]interface{}{
					"parameters": map[string]interface{}{
						"in":       "path",
						"requried": true,
						"schema": map[string]interface{}{
							"$ref": "#/components/schemas/ref1",
						},
					},
					"tags": []interface{}{
						"tag5",
					},
				},
				"post": map[string]interface{}{
					"tags": []interface{}{
						"tag5",
					},
				},
			},
			"/path2": map[string]interface{}{
				"delete": map[string]interface{}{
					"tags": []interface{}{
						"tag2",
					},
				},
			},
			"/path32": map[string]interface{}{
				"get": map[string]interface{}{
					"tags": []interface{}{
						"tag3",
					},
				},
			},
			"/path4": map[string]interface{}{
				"post": map[string]interface{}{
					"tags": []interface{}{
						"tag4",
					},
				},
			},
			"/path31": map[string]interface{}{
				"delete": map[string]interface{}{
					"tags": []interface{}{
						"tag3",
					},
				},
				"get": map[string]interface{}{
					"tags": []interface{}{
						"tag3",
					},
				},
				"patch": map[string]interface{}{
					"responses": map[string]interface{}{
						"500": map[string]interface{}{
							"description": "Internal Server Error",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/ref3",
									},
								},
							},
						},
					},
					"tags": []interface{}{
						"tag3",
					},
				},
			},
		},
		"components": map[string]interface{}{
			"schemas": map[string]interface{}{
				"ref1": map[string]interface{}{
					"type": "object",
				},
				"ref2": map[string]interface{}{
					"type":  "array",
					"title": "ref2Name",
				},
				"ref3": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"tags": []interface{}{
			map[string]interface{}{
				"name": "tag4",
			},
			map[string]interface{}{
				"name": "tag1",
			},
			map[string]interface{}{
				"name": "tag3",
				"externalDocs": map[string]string{
					"url": "https://blah-blah.ru/endpoint",
				},
			},
			map[string]interface{}{
				"name": "tag5",
			},
			map[string]interface{}{
				"name": "tag2",
			},
		},
	}
}

func TestUpdateSwagger(t *testing.T) {
	cases := []struct {
		in   updateSwaggerArguments
		want map[string]interface{}
	}{
		{
			in: updateSwaggerArguments{
				swagger:      map[string]interface{}{},
				pathesToKeep: []S.PathMethod{},
			},
			want: map[string]interface{}{},
		},
		{
			in: updateSwaggerArguments{
				swagger: testSwaggerFactory(),
				pathesToKeep: []S.PathMethod{
					{
						Path:   "/path5",
						Method: "delete",
					},
				},
			},
			want: map[string]interface{}{
				"paths": map[string]interface{}{
					"/path5": map[string]interface{}{
						"delete": map[string]interface{}{
							"tags": []interface{}{
								"tag5",
							},
						},
					},
				},
				"components": map[string]interface{}{
					"schemas": map[string]interface{}{
						"ref1": map[string]interface{}{
							"type": "object",
						},
					},
				},
				"tags": []map[string]interface{}{
					{
						"name": "tag5",
					},
				},
			},
		},
		{
			in: updateSwaggerArguments{
				swagger: testSwaggerFactory(),
				pathesToKeep: []S.PathMethod{
					{
						Path:   "/path5",
						Method: "delete",
					},
					{
						Path:   "/path31",
						Method: "patch",
					},
				},
			},
			want: map[string]interface{}{
				"paths": map[string]interface{}{
					"/path31": map[string]interface{}{
						"patch": map[string]interface{}{
							"tags": []interface{}{
								"tag3",
							},
						},
					},
					"/path5": map[string]interface{}{
						"delete": map[string]interface{}{
							"tags": []interface{}{
								"tag5",
							},
						},
					},
				},
				"components": map[string]interface{}{
					"schemas": map[string]interface{}{
						"ref1": map[string]interface{}{
							"type": "object",
						},
						"ref3": map[string]interface{}{
							"type": "string",
						},
					},
				},
				"tags": []map[string]interface{}{
					{
						"name": "tag3",
						"externalDocs": map[string]string{
							"url": "https://blah-blah.ru/endpoint",
						},
					},
					{
						"name": "tag5",
					},
				},
			},
		},
	}

	for currentCaseIdx, currentCase := range cases {
		got, err := UpdateSwagger(currentCase.in.swagger, currentCase.in.pathesToKeep)

		gotLen := len(got)
		wantLen := len(currentCase.want)
		if err != nil && wantLen != 0 {
			t.Errorf(
				"Testing UpdateSwagger:\nexpected expected positive case, but got error: \n%+v",
				err,
			)
		}

		if wantLen != gotLen {
			t.Errorf(
				"Testing UpdateSwagger:\nexpected map has different amount of keys from what we got\nexpected: %d\ngot: %d\n",
				wantLen,
				gotLen,
			)
		}

		if gotLen == 0 {
			continue
		}

		gotTagsI, ok := got["tags"]
		if !ok {
			t.Error("Testing UpdateSwagger:\nexpected to get map with tags, but something gone wrong")
		}
		gotTags, ok := gotTagsI.([]map[string]interface{})
		if !ok {
			t.Error("Testing UpdateSwagger:\nexpected to get map with tags, but something gone wrong")
		}

		wantTagsI := currentCase.want["tags"]
		wantTags, ok := wantTagsI.([]map[string]interface{})
		if !ok {
			t.Error("Testing UpdateSwagger:\nexpected to get map with tags, but something gone wrong")
		}

		gotTagsLen := len(gotTags)
		wantTagsLen := len(wantTags)
		if wantTagsLen != gotTagsLen {
			t.Errorf(
				"Testing UpdateSwagger:\ninputed pathes: %+v\nexpected tags length: %v\ngot tags length %v",
				currentCase.in.pathesToKeep,
				wantTagsLen,
				gotTagsLen,
			)
		}
		for tagIdx, wantTag := range wantTags {
			wantTagName := wantTag["name"]
			gotTagName := gotTags[tagIdx]["name"]
			if wantTagName != gotTagName {
				t.Errorf(
					"Testing UpdateSwagger:\ninputed pathes: %+v\nincorrect name of tags[%d]:\nexpected: %s\ngot: %s",
					currentCase.in.pathesToKeep,
					tagIdx,
					wantTagName,
					gotTagName,
				)
			}
		}

		wantComponentsI := currentCase.want["components"]
		if wantComponents, ok := wantComponentsI.(map[string]interface{}); ok {
			wantSchemasI := wantComponents["schemas"]
			wantSchemas, ok := wantSchemasI.(map[string]interface{})
			if !ok {
				t.Errorf("Something went wrong with case %d: cant extract schemas", currentCaseIdx)
			}

			gotComponentsI, ok := got["components"]
			if !ok {
				t.Errorf("Expected to have components in result's swagger, but got nogthing in case %d", currentCaseIdx)
			}
			gotComponents, ok := gotComponentsI.(map[string]interface{})
			if !ok {
				t.Errorf("Components in result's swagger has wrong type in case %d", currentCaseIdx)
			}

			gotSchemasI, ok := gotComponents["schemas"]
			if !ok {
				t.Errorf("Expected to have schemas in result's swagger, but got nothing in case %d", currentCaseIdx)
			}
			gotSchemas, ok := gotSchemasI.(map[string]interface{})
			if !ok {
				t.Errorf("Schemas in result's swagger has wrong type in case %d", currentCaseIdx)
			}
			wantSchemasLen := len(wantSchemas)
			gotSchemasLen := len(gotSchemas)
			if gotSchemasLen != wantSchemasLen {
				t.Errorf(
					"Testing UpdateSwagger: case %d\nexpected schemas length: %v\ngot schemas length %v",
					currentCaseIdx,
					wantSchemasLen,
					gotSchemasLen,
				)
			}

			for wantSchemaName := range wantSchemas {
				_, hasSchemaName := gotSchemas[wantSchemaName]
				if !hasSchemaName {
					t.Errorf(
						"Testing UpdateSwagger: case %d\nexpected result to have schema %s, but it doesn't",
						currentCaseIdx,
						wantSchemaName,
					)
				}
			}
		}
	}
}
