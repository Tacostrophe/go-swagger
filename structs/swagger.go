package structs

type PathMethod struct {
	Path   string
	Method string
}

// type SwaggerInfo struct {
// 	title          string
// 	description    string
// 	termsOfService string
// contact        Contact
// license        License
// 	version string
// }

type Swagger struct {
	Paths map[string]map[string]interface{} `json: "paths"`
}
