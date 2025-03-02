package usecases

type (
	SwaggerUsecase interface {
		Init(string) error
		GetFilteredPathes(string) []PathMethod
	}
)
