package bear

import "runtime"

// Template is an error template
type Template struct {
	opts []ErrOption
}

// NewTemplate creates a new error template
func NewTemplate(opts ...ErrOption) Template {
	return Template{
		opts: opts,
	}
}

// NewTemplate creates a new template with the current template as a base
func (t *Template) NewTemplate(opts ...ErrOption) Template {
	return Template{
		opts: append(t.opts, opts...),
	}
}

// Union performs unions the given template with the current template
func (t *Template) Union(template *Template) {
	t.opts = append(t.opts, template.opts...)
}

// New creates a new error from the template
func (t *Template) New(opts ...ErrOption) *Error {
	err := New(append(t.opts, opts...)...)

	// reset the stack trace since were calling new from inside the package
	err.stack = getStackTrace(2)

	return err
}

// getStackTrage get's all the stack frames the lead upto an error being called
func getStackTrace(initialSkip int) []stackFrame {
	i := initialSkip
	var frames []stackFrame
	for {
		_, filename, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		frames = append(frames, stackFrame{
			filename: filename,
			line:     line,
		})

		i++
	}

	return frames
}
