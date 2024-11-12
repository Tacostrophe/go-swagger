package init_context

import (
	"testing"

	S "github.com/Tacostrophe/go-swagger/structs"
)

func TestInitContext(t *testing.T) {
	cases := []struct {
		in   []string
		want struct {
			value    S.Context
			hasError bool
		}
	}{
		{
			in: []string{
				"path/to/go",
			},
			want: struct {
				value    S.Context
				hasError bool
			}{
				value:    S.Context{},
				hasError: true,
			},
		},
		{
			in: []string{
				"path/to/go",
				"path/to/swagger",
			},
			want: struct {
				value    S.Context
				hasError bool
			}{
				value:    S.Context{},
				hasError: true,
			},
		},
		{
			in: []string{
				"path/to/go",
				"path/to/swagger.json",
			},
			want: struct {
				value    S.Context
				hasError bool
			}{
				value: S.Context{
					IncomeSwaggerPath: "path/to/swagger.json",
				},
				hasError: false,
			},
		},
	}

	for _, currentCase := range cases {
		got, err := InitContext(currentCase.in)

		if got != currentCase.want.value {
			t.Errorf("expected InitContext(%+v)\n  to return: %+v\n  but got: %+v", currentCase.in, currentCase.want.value, got)
		}
		if currentCase.want.hasError && err == nil {
			t.Errorf("expected InitContext(%+v)\n  to return error\n  but got nil", currentCase.in)
		}
		if !currentCase.want.hasError && err != nil {
			t.Errorf("expected InitContext(%+v)\n  NOT to return error\n  but got error: %+v", currentCase.in, err)
		}
	}
}
