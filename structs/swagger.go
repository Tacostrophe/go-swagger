package structs

type (
	Tags struct {
		Tag []PathesTagged
	}

	PathesTagged struct {
		tagName string
		pathes  []PathMethod
	}

	PathMethod struct {
		Path     string
		Method   string
		FirstTag string
	}

	Swagger struct {
		Paths map[string]map[string]interface{} `json: "paths"`
	}
)

func NewPathesTagged(pathes []PathMethod) {
}

// func (p *Pathes) Get() []PathMethod {
// 	return p.pathes
// }

// func (p *Pathes) Filter(filters string) []PathMethod {
// 	if filters == "" {
// 		return p.pathes
// 	}
// }
