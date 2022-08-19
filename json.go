package bear

import (
	"sort"

	"github.com/bjatkin/bear/pkg/metrics"
)

// jsonError mirriors the Error type but it's fields are exported so it can be json marshled
type jsonError struct {
	ID       *string                `json:"id,omitempty"`
	Parents  []jsonError            `json:"parents,omitempty"`
	ErrType  *ErrType               `json:"errType,omitempty"`
	Tags     map[string]interface{} `json:"tags,omitempty"`
	Labels   []string               `json:"labels,omitempty"`
	Metrics  []*metrics.Metric      `json:"metrics,omitempty"`
	Fmetrics []*metrics.FMetric     `json:"fmetrics,omitempty"`
	Msg      *string                `json:"msg,omitempty"`
	Code     *int                   `json:"code,omitempty"`
	ExitCode *int                   `json:"exitCode,omitempty"`
	Stack    []string               `json:"stack,omitempty"`
}

// newJSONError creates a new jsonError from an Error
func newJSONError(e *Error) jsonError {
	err := jsonError{
		ID:       &e.id,
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
		err.Msg = nil
	}

	if e.noID {
		err.ID = nil
	}

	if !e.noParents {
		err.Parents = buildParents(e.parents, e.noStack, e.noMsg, e.noID)
	}

	if !e.noStack {
		for _, frame := range e.stack {
			err.Stack = append(err.Stack, frame.String())
		}
	}

	return err
}

// buildParents builds an array of jsonErrors from a slice of parent errors
func buildParents(parents []error, noStack, noMsg, noID bool) []jsonError {
	var jsonParents []jsonError
	for _, parent := range parents {
		berr, _ := AsBerr(parent)
		jerr := newJSONError(berr)
		if noStack {
			jerr.Stack = nil
		}
		if noMsg {
			jerr.Msg = nil
		}
		if noID {
			jerr.ID = nil
		}

		jsonParents = append(jsonParents, jerr)
	}

	return jsonParents
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
