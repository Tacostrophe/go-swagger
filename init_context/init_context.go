package init_context

import (
	"errors"
	"strings"

	S "github.com/Tacostrophe/go-swagger/structs"
)

func InitContext(args []string) (S.Context, error) {
	if len(args) < 2 {
		return S.Context{}, errors.New("path/to/json as first argument is required")
	}
	incomeSwaggerPath := args[1]
	if !strings.HasSuffix(incomeSwaggerPath, ".json") {
		return S.Context{}, errors.New("file must be a json")
	}

	ctx := S.Context{
		IncomeSwaggerPath: incomeSwaggerPath,
	}

	return ctx, nil
}
