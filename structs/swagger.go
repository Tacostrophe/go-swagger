package structs

type PathMethod struct {
	Path     string
	Method   string
	FirstTag string
}

type Swagger struct {
	Paths map[string]map[string]interface{} `json: "paths"`
}
