package bear

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/bjatkin/bear/pkg/metrics"
)

var (
	panicErr = NewType("Panic Error")
)

// Error is a custom bear error
type Error struct {
	parents  []error
	errType  *ErrType
	tags     map[string]interface{}
	labels   map[string]struct{}
	metrics  []*metrics.Metric
	fmetrics []*metrics.FMetric
	msg      *string
	code     *int
	exitCode *int
	stack    []stackFrame

	// fmt settings
	prettyPrint bool
	noStack     bool
	noParents   bool
	noMsg       bool

	// panic settings
	stdErr io.Writer
}

// stackFrame is a stack frame in the codes execution
type stackFrame struct {
	filename string
	line     int
}

func (f stackFrame) String() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("failed to get current working dir: " + err.Error())
	}
	return fmt.Sprintf("%s:%d", strings.TrimPrefix(f.filename, dir+"/"), f.line)
}

// New creates a new bear.Error
func New(opts ...ErrOption) *Error {
	e := &Error{
		stack:  getStackTrace(2),
		stdErr: os.Stderr,
	}
	for _, opt := range opts {
		opt(e)
	}

	return e
}

// Wrap creates a new bear.Error with parent err
func Wrap(err error, opts ...ErrOption) *Error {
	e := New(append(opts, WithParent(err))...)

	// reset the stack trace since were calling new from inside the package
	e.stack = getStackTrace(2)

	return e
}

// ErrOption adds optional info to an error
type ErrOption func(*Error)

// WithParent adds a parent error
func WithParent(err error) ErrOption {
	return func(e *Error) {
		e.parents = append(e.parents, err)
	}
}

// WithCode adds the error code to the error
func WithCode(code int) ErrOption {
	return func(e *Error) {
		e.code = &code
	}
}

// WithExitCode sets the exit code for the error
func WithExitCode(exitCode int) ErrOption {
	return func(e *Error) {
		e.exitCode = &exitCode
	}
}

// WithTag adds tag information to an error
func WithTag(name string, value interface{}) ErrOption {
	return func(e *Error) {
		if e.tags == nil {
			e.tags = make(map[string]interface{})
		}
		e.tags[name] = value
	}
}

// WithLabels adds one or more labels to the error
func WithLabels(names ...string) ErrOption {
	return func(e *Error) {
		if e.labels == nil {
			e.labels = make(map[string]struct{})
		}
		for _, name := range names {
			e.labels[name] = struct{}{}
		}
	}
}

// WithMetrics adds new metrics to the error
func WithMetrics(metrics ...*metrics.Metric) ErrOption {
	return func(e *Error) {
		e.metrics = append(e.metrics, metrics...)
	}
}

// WithFMetric adds new fmetrics to the error
func WithFMetric(metrics ...*metrics.FMetric) ErrOption {
	return func(e *Error) {
		e.fmetrics = append(e.fmetrics, metrics...)
	}
}

// ErrType gives the error a general error class
type ErrType string

// NewType creates a new error type, these should be broad categories of errors rather than specific
func NewType(name string) ErrType {
	return ErrType(name)
}

// WithType adds an error type to the error
func WithErrType(t ErrType) ErrOption {
	return func(e *Error) {
		e.errType = &t
	}
}

// WithMsg adds the message to the error this should only be used rarely
func WithMsg(message string) ErrOption {
	return func(e *Error) {
		e.msg = &message
	}
}

// WithStdErr changes the stdErr output to the provided buffer
func WithStdErr(buf io.Writer) ErrOption {
	return func(e *Error) {
		e.stdErr = buf
	}
}

// FmtPrettyPrint tells the error to format its error
func FmtPrettyPrint(on bool) ErrOption {
	return func(e *Error) {
		e.prettyPrint = on
	}
}

// FmtNoStack turns off the stack trace for the error
func FmtNoStack(on bool) ErrOption {
	return func(e *Error) {
		e.noStack = on
	}
}

// FmtNoParents turns off the parents for error
func FmtNoParents(on bool) ErrOption {
	return func(e *Error) {
		e.noParents = on
	}
}

// FmtNoMsg turns off the message for error
func FmtNoMsg(on bool) ErrOption {
	return func(e *Error) {
		e.noMsg = on
	}
}

// Error implements the error interface
func (e *Error) Error() string {
	// create a public type for json marshaling
	// is there a better way to do this?
	public := struct {
		Parents  []json.RawMessage      `json:"parents,omitempty"`
		ErrType  *ErrType               `json:"errType,omitempty"`
		Tags     map[string]interface{} `json:"tags,omitempty"`
		Labels   []string               `json:"labels,omitempty"`
		Metrics  []*metrics.Metric      `json:"metrics,omitempty"`
		Fmetrics []*metrics.FMetric     `json:"fmetrics,omitempty"`
		Msg      *string                `json:"msg,omitempty"`
		Code     *int                   `json:"code,omitempty"`
		ExitCode *int                   `json:"exitCode,omitempty"`
		Stack    []string               `json:"stack,omitempty"`
	}{
		ErrType:  e.errType,
		Tags:     e.tags,
		Labels:   mapToArray(e.labels),
		Metrics:  e.metrics,
		Fmetrics: e.fmetrics,
		Msg:      e.msg,
		Code:     e.code,
		ExitCode: e.exitCode,
	}

	if e.noMsg {
		public.Msg = nil
	}

	if !e.noParents {
		for _, parent := range e.parents {
			public.Parents = append(public.Parents, json.RawMessage([]byte(parent.Error())))
		}
	}

	if !e.noStack {
		for _, frame := range e.stack {
			public.Stack = append(public.Stack, frame.String())
		}
	}

	// print data as json
	var raw []byte
	var err error
	if e.prettyPrint {
		raw, err = json.MarshalIndent(public, "", "  ")
		if err != nil {
			panic("failed to marshal indent error: " + err.Error())
		}
	} else {
		raw, err = json.Marshal(public)
		if err != nil {
			panic("failed to marshal error: " + err.Error())
		}
	}

	return string(raw)
}

// mapToArray converts the map to an array
func mapToArray(m map[string]struct{}) []string {
	var slice []string
	for m := range m {
		slice = append(slice, m)
	}

	// sort the slice for consistent label order
	sort.Strings(slice)

	return slice
}

// Panic will conver the error into a panic, if print is true the error will be printed to stdErr
func (e *Error) Panic(print bool) {
	if print {
		_, _ = fmt.Fprintln(e.stdErr, e.Error())
	}

	switch {
	case e.exitCode != nil && *e.exitCode != 0:
		panic(e.exitCode)
	case e.code != nil && *e.code != 0:
		panic(e.exitCode)
	default:
		panic(1)
	}
}

// Add adds new options to the error
func (e *Error) Add(opts ...ErrOption) *Error {
	for _, opt := range opts {
		opt(e)
	}

	return e
}

// HasLabel returns true if the lable has been set on the error
func (e *Error) HasLabel(label string) bool {
	if e.labels == nil {
		return false
	}

	_, ok := e.labels[label]
	return ok
}

// HasTag returns true if the tag has been set on the error
func (e *Error) HasTag(tag string) bool {
	if e.tags == nil {
		return false
	}

	_, ok := e.tags[tag]
	return ok
}

// GetTag returns true and the value of the tag if it has been set on the error
// otherwise nil, and false are returned
func (e *Error) GetTag(tag string) (interface{}, bool) {
	if e.tags == nil {
		return nil, false
	}

	i, ok := e.tags[tag]
	return i, ok
}

// WrapPanic should be used as a defer function, it will catch any panics inside the function as wrap them as a parent error
// the provided opts are used to create the new panic error
func (e *Error) WrapPanic(opts ...ErrOption) {
	if err := recover(); err != nil {
		parent := New(opts...).Add(WithErrType(panicErr))

		switch v := err.(type) {
		// all the int types
		case int:
			parent.Add(WithExitCode(v))
		case int8:
			parent.Add(WithExitCode(int(v)))
		case int16:
			parent.Add(WithExitCode(int(v)))
		case int32:
			parent.Add(WithExitCode(int(v)))
		case int64:
			parent.Add(WithExitCode(int(v)))

		// all the uint types
		case uint:
			parent.Add(WithExitCode(int(v)))
		case uint8:
			parent.Add(WithExitCode(int(v)))
		case uint16:
			parent.Add(WithExitCode(int(v)))
		case uint32:
			parent.Add(WithExitCode(int(v)))
		case uint64:
			parent.Add(WithExitCode(int(v)))

		case string:
			parent.Add(WithMsg(v))

		default:
			parent.Add(WithTag("value", v), WithTag("type", fmt.Sprintf("%T", v)))
		}

		e.Add(WithParent(parent))
	}
}
