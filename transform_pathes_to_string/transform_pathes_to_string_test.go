package transform_pathes_to_string

import (
	"testing"

	S "github.com/Tacostrophe/go-swagger/structs"
)

func TestTransformPathesToString(t *testing.T) {
	cases := []struct {
		in   []S.PathMethod
		want string
	}{
		{
			in: []S.PathMethod{
				{
					Path:   "test/path/1",
					Method: "get",
				},
			},
			want: "  0. test/path/1 GET",
		},
		{
			in: []S.PathMethod{
				{
					Path:   "test/path/1",
					Method: "delete",
				},
				{
					Path:   "test/path/1",
					Method: "get",
				},
				{
					Path:   "test/path/2",
					Method: "get",
				},
			},
			want: "  0. test/path/1 DELETE\n  1.             GET\n  2. test/path/2 GET",
		},
	}

	for _, currentCase := range cases {
		got, _ := TransformPathesToString(currentCase.in)

		if currentCase.want != got {
			t.Errorf(
				"Testing TransformPthaesToString:\ninput: %+v\nexpected: -\ngot: +\n- \"%s\"\n+ \"%s\"",
				currentCase.in,
				currentCase.want,
				got,
			)
		}
	}
}
