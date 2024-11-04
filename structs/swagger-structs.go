package structs

type PathMethod struct {
	Path   string
	Method string
}

type SwaggerInfo struct {
	title          string
	description    string
	termsOfService string
	// contact        Contact
	// license        License
	version string
}

type Path struct {
}

type Swagger struct {
	openapi string
	info    SwaggerInfo
	pathes  map[string]map[string]Path
}
