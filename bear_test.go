package bear

import (
	"bytes"
	"strings"
	"testing"
)

// TODO: test stack trace in it's own function

func TestNew(t *testing.T) {
	defaultOpts := []ErrOption{FmtNoStack(true)}
	parent := NewTemplate(defaultOpts...)

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
			"with parent error",
			args{
				opts: append(defaultOpts, WithParent(parent.New(WithCode(1)))),
			},
			nil,
			`{"parents":[{"code":1}]}`,
		},
		{
			"with exit code",
			args{
				opts: append(defaultOpts, WithExitCode(400)),
			},
			nil,
			`{"exitCode":400}`,
		},
		{
			"with tags",
			args{
				opts: append(defaultOpts, WithTag("test", 1), WithTag("final", true), WithTag("hello", "world")),
			},
			nil,
			`{"tags":{"final":true,"hello":"world","test":1}}`,
		},
		{
			"with labels",
			args{
				opts: append(defaultOpts, WithLabel("test"), WithLabel("default")),
			},
			nil,
			`{"labels":["test","default"]}`,
		},
		{
			"with error type",
			args{
				opts: append(defaultOpts, WithErrType(NewType("test error"))),
			},
			nil,
			`{"errType":"test error"}`,
		},
		{
			"with message",
			args{
				opts: append(defaultOpts, WithMsg("this is a test message")),
			},
			nil,
			`{"msg":"this is a test message"}`,
		},
		{
			"with pretty print",
			args{
				opts: append(defaultOpts,
					WithParent(parent.New(WithCode(2))),
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
			"with not parents",
			args{
				opts: append(defaultOpts, WithParent(parent.New(WithCode(3))), FmtNoParents(true)),
			},
			nil,
			`{}`,
		},
		{
			"with stack",
			args{},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.opts...)
			// run any final setup before we check the output
			if tt.setup != nil {
				tt.setup(got)
			}

			if got.Error() != tt.want {
				t.Errorf("New() error string = \n%s, want \n%s", got.Error(), tt.want)
			}
		})
	}
}

func TestTemplate_New(t *testing.T) {
	type fields struct {
		opts []ErrOption
	}
	type args struct {
		opts []ErrOption
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"empty template",
			fields{
				opts: []ErrOption{},
			},
			args{
				opts: []ErrOption{WithCode(1), FmtNoStack(true)},
			},
			`{"code":1}`,
		},
		{
			"no stack template",
			fields{
				opts: []ErrOption{FmtNoStack(true)},
			},
			args{
				opts: []ErrOption{WithCode(1)},
			},
			`{"code":1}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewTemplate(tt.fields.opts...)
			if got := tr.New(tt.args.opts...).Error(); got != tt.want {
				t.Errorf("Template.New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_WrapPanic(t *testing.T) {
	tests := []struct {
		name      string
		panicWith interface{}
		want      string
	}{
		{
			"panic with error code",
			1,
			`{"parents":[{"errType":"Panic Error","exitCode":1}]}`,
		},
		{
			"panic with uint64",
			uint64(2),
			`{"parents":[{"errType":"Panic Error","exitCode":2}]}`,
		},
		{
			"panic with string",
			"panicing",
			`{"parents":[{"errType":"Panic Error","msg":"panicing"}]}`,
		},
		{
			"panic with bool",
			true,
			`{"parents":[{"errType":"Panic Error","tags":{"type":"bool","value":true}}]}`,
		},
		{
			"panic with struct",
			struct {
				test bool
				str  string
			}{true, "test"},
			`{"parents":[{"errType":"Panic Error","tags":{"type":"struct { test bool; str string }","value":{}}}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := func() (e error) {
				e = New(FmtNoStack(true))
				defer (e.(*Error)).WrapPanic(FmtNoStack(true))

				panic(tt.panicWith)
			}()

			if got := err.Error(); got != tt.want {
				t.Fatalf("WrapPanic() got \n%s, want \n%s", got, tt.want)
			}
		})
	}
}

func TestError_Panic(t *testing.T) {
	type fields struct {
		opts []ErrOption
	}
	type args struct {
		print bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"don't print",
			fields{
				opts: []ErrOption{WithExitCode(1)},
			},
			args{
				print: false,
			},
			"",
		},
		{
			"print error",
			fields{
				opts: []ErrOption{WithCode(1), FmtNoStack(true)},
			},
			args{
				print: true,
			},
			`{"code":1}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdErr := &bytes.Buffer{}
			e := New(tt.fields.opts...).Add(WithStdErr(stdErr))

			// catch the expected panic
			defer func() {
				if err := recover(); err != nil {
					// success, nothing to do here
				} else {
					t.Errorf("Panic() calling defered function code did not panic")
				}

				got := strings.Trim(stdErr.String(), "\n")
				if got != tt.want {
					t.Errorf("Panic() stdErr was %s, wanted %s", got, tt.want)
				}
			}()

			e.Panic(tt.args.print)

			t.Error("Panic() code did not panic")
		})
	}
}
