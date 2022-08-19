package bear

import (
	"bytes"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/bjatkin/bear/pkg/metrics"
)

func TestNew(t *testing.T) {
	// default options to make testing easier
	defaultOpts := []ErrOption{FmtNoStack(true), FmtNoID(true)}
	hexReg := regexp.MustCompile(`[0-9a-f]{64}`)

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
				opts: append(defaultOpts, WithParent(New(WithCode(1)))),
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
				opts: append(defaultOpts, WithLabels("test", "default")),
			},
			nil,
			`{"labels":["default","test"]}`,
		},
		{
			"with metrics",
			args{
				opts: append(defaultOpts, WithMetrics(metrics.NewMetric("test"), metrics.NewMetric("again"))),
			},
			nil,
			`{"metrics":[{"name":"test","value":0},{"name":"again","value":0}]}`,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.opts...)

			// check to make sure the id was set to a valid hex string
			if !hexReg.MatchString(got.GetID()) {
				t.Errorf("New() error id was not set correctly")
			}

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

func TestError_WrapPanic(t *testing.T) {
	tests := []struct {
		name      string
		panicWith interface{}
		want      string
		wantErr   bool
	}{
		{
			"panic with error code",
			1,
			`{"parents":[{"errType":"Panic Error","exitCode":1}]}`,
			true,
		},
		{
			"panic with uint64",
			uint64(2),
			`{"parents":[{"errType":"Panic Error","exitCode":2}]}`,
			true,
		},
		{
			"panic with string",
			"panicing",
			`{"parents":[{"errType":"Panic Error","msg":"panicing"}]}`,
			true,
		},
		{
			"panic with bool",
			true,
			`{"parents":[{"errType":"Panic Error","tags":{"type":"bool","value":true}}]}`,
			true,
		},
		{
			"panic with struct",
			struct {
				test bool
				str  string
			}{true, "test"},
			`{"parents":[{"errType":"Panic Error","tags":{"type":"struct { test bool; str string }","value":{}}}]}`,
			true,
		},
		{
			"no panic",
			nil,
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := func() (e error) {
				e = New(FmtNoStack(true), FmtNoID(true))
				defer (e.(*Error)).WrapPanic(FmtNoStack(true))

				if tt.wantErr {
					panic(tt.panicWith)
				}

				return nil
			}()

			if (err != nil) != tt.wantErr {
				t.Fatalf("WrapPanic() got err %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				// no need to do further checks on the error if we didn't even want it
				return
			}

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
				opts: []ErrOption{WithCode(1), FmtNoStack(true), FmtNoID(true)},
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

func TestError_Add(t *testing.T) {
	type args struct {
		opts []ErrOption
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"add nothing",
			args{},
			`{}`,
		},
		{
			"add code and labels",
			args{
				opts: []ErrOption{WithCode(1), WithLabels("test", "success")},
			},
			`{"labels":["success","test"],"code":1}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(FmtNoStack(true), FmtNoID(true)).Add(tt.args.opts...)

			got := e.Add(tt.args.opts...)

			gotErr := got.Error()
			if gotErr != tt.want {
				t.Errorf("Error.Add() error string was \n%v, want \n%v", gotErr, tt.want)
			}
		})
	}
}

func TestError_HasLabel(t *testing.T) {
	type fields struct {
		labels map[string]struct{}
	}
	type args struct {
		label string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"nil labels",
			fields{},
			args{
				label: "nil",
			},
			false,
		},
		{
			"missing label",
			fields{
				labels: map[string]struct{}{"success": {}},
			},
			args{
				label: "missing",
			},
			false,
		},
		{
			"found label",
			fields{
				labels: map[string]struct{}{"success": {}},
			},
			args{
				label: "success",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				labels: tt.fields.labels,
			}

			if got := e.HasLabel(tt.args.label); got != tt.want {
				t.Errorf("Error.HasLabel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_HasTag(t *testing.T) {
	type fields struct {
		tags map[string]interface{}
	}
	type args struct {
		tag string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"nil tags",
			fields{},
			args{
				tag: "nil",
			},
			false,
		},
		{
			"missing tag",
			fields{
				tags: map[string]interface{}{"success": false},
			},
			args{
				tag: "missing",
			},
			false,
		},
		{
			"found tag",
			fields{
				tags: map[string]interface{}{"success": true},
			},
			args{
				tag: "success",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				tags: tt.fields.tags,
			}
			if got := e.HasTag(tt.args.tag); got != tt.want {
				t.Errorf("Error.HasTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_GetTag(t *testing.T) {
	type fields struct {
		tags map[string]interface{}
	}
	type args struct {
		tag string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
		wantOK bool
	}{
		{
			"nil tags",
			fields{},
			args{
				tag: "nil",
			},
			nil,
			false,
		},
		{
			"missing tag",
			fields{
				map[string]interface{}{"success": false},
			},
			args{
				tag: "missing",
			},
			nil,
			false,
		},
		{
			"found tag",
			fields{
				map[string]interface{}{"success": true},
			},
			args{
				tag: "success",
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				tags: tt.fields.tags,
			}

			got, gotOK := e.GetTag(tt.args.tag)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.GetTag() got = %v, want %v", got, tt.want)
			}
			if gotOK != tt.wantOK {
				t.Errorf("Error.GetTag() gotOK = %v, wantOK %v", gotOK, tt.wantOK)
			}
		})
	}
}
