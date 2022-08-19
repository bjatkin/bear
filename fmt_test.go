package bear

import (
	"testing"
)

func TestNewFmt(t *testing.T) {
	// default options to make testing easier
	defaultOpts := []ErrOption{FmtNoStack(true), FmtNoID(true)}

	type args struct {
		opts []ErrOption
	}
	tests := []struct {
		name  string
		args  args
		setup func(e *Error)
		want  string
	}{
		{
			"with pretty print",
			args{
				opts: append(defaultOpts,
					WithParent(New(WithCode(2))),
					WithExitCode(404),
					FmtPrettyPrint(true),
				),
			},
			nil,
			`{
  "parents": [
    {
      "code": 2
    }
  ],
  "exitCode": 404
}`,
		},
		{
			"with no parents",
			args{
				opts: append(defaultOpts, WithParent(New(WithCode(3))), FmtNoParents(true)),
			},
			nil,
			`{}`,
		},
		{
			"with stack",
			args{
				opts: []ErrOption{FmtNoID(true)},
			},
			func(e *Error) {
				e.stack = []stackFrame{
					{filename: "test.go", line: 100},
					{filename: "test2.go", line: 50},
					{filename: "final.go", line: 1},
				}
			},
			`{"stack":["test.go:100","test2.go:50","final.go:1"]}`,
		},
		{
			"with no message",
			args{
				opts: append(defaultOpts, WithMsg("this is a test message"), FmtNoMsg(true)),
			},
			nil,
			`{}`,
		},
		{
			"with id",
			args{
				opts: []ErrOption{FmtNoStack(true)},
			},
			func(e *Error) {
				e.id = "test"
			},
			`{"id":"test"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.opts...)
			// run any final setup before we check the output
			if tt.setup != nil {
				tt.setup(got)
			}

			gotErr := got.Error()
			if gotErr != tt.want {
				t.Errorf("New() error string was \n'%s', want \n'%s'", gotErr, tt.want)
			}
		})
	}
}
