package bear

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bjatkin/bear/pkg/metrics"
)

var (
	PanicErr = NewType("Panic Error")
)

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

// Error is a custom bear error
type Error struct {
	id       string
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
	noID        bool

	// panic settings
	stdErr io.Writer
}

// New creates a new bear.Error
func New(opts ...ErrOption) *Error {
	e := &Error{
		id:     newRandomID(),
		stack:  getStackTrace(2),
		stdErr: os.Stderr,
	}
	for _, opt := range opts {
		opt(e)
	}

	return e
}

func newRandomID() string {
	randSource := rand.NewSource(time.Now().UnixNano())
	idLen := 64
	ref := "abcdef0123456789"

	var id string
	for i := 0; i < idLen; i++ {
		id += string(ref[randSource.Int63()%int64(len(ref))])
	}

	return id
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

// Error implements the error interface
func (e *Error) Error() string {
	// create a jsonError for json marshaling
	public := newJSONError(e)

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

// GetID returns the ID of the error
func (e *Error) GetID() string {
	return e.id
}

// WrapPanic should be used as a defer function, it will catch any panics inside the function as wrap them as a parent error
// the provided opts are used to create the new panic error
func (e *Error) WrapPanic(opts ...ErrOption) {
	if err := recover(); err != nil {
		parent := New(opts...).Add(WithErrType(PanicErr))

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

// IsBerr returns true if the given error is a bear error
func IsBerr(e error) bool {
	if _, ok := e.(*Error); ok {
		return true
	}
	return false
}

// AsBerr returns the error and true if the given error is a bear error
// otherwise it will convert the error into a bear error and return false
func AsBerr(e error) (*Error, bool) {
	if berr, ok := e.(*Error); ok {
		return berr, ok
	}

	return New(WithMsg(e.Error())), true
}

// Is returns true if the error is of the same type as the error provided
// If e is a bear.Error the errType field will be checked
// otherwise the e.Error() method will be check to see if it matches the error type
func Is(e error, t ErrType) bool {
	if berr, ok := AsBerr(e); ok {
		return berr.errType != nil && *berr.errType == t
	}

	return e.Error() == string(t)
}
