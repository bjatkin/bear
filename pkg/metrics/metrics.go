package metrics

import (
	"encoding/json"
	"fmt"
)

// Metric is a tracker that usese an int to track a value
type Metric struct {
	name   string
	metric int
}

// NewMetric creates a new metric
func NewMetric(name string) *Metric {
	return &Metric{
		name: name,
	}
}

// Incr adds one to the metric
func (m *Metric) Incr() {
	m.metric++
}

// Decr subtrcts one from the metric
func (m *Metric) Decr() {
	m.metric--
}

// Add adds an abitrary value to the underlying metric
func (m *Metric) Add(add int) {
	m.metric += add
}

// MarshalJSON implements the marshaler interface
func (m *Metric) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}{
		m.name,
		m.metric,
	})
}

// String implements the stringer interface
func (m *Metric) String() string {
	return fmt.Sprintf("[%s] %d", m.name, m.metric)
}

// GetName returns the name of the metric
func (m *Metric) GetName() string {
	return m.name
}

// GetValue returns the value of the metric
func (m *Metric) GetValue() int {
	return m.metric
}

// AddMetrics adds metric values together and returns the result
func AddMetrics(metrics ...Metric) int {
	var total int
	for _, m := range metrics {
		total += m.metric
	}
	return total
}

// FilterMetrics takes a list of metrics and returns only those that satisfy the filter metric
func FilterMetrics(metrics []Metric, filter func(Metric) bool) []Metric {
	if filter == nil {
		return metrics
	}

	var filtered []Metric
	for _, metric := range metrics {
		if filter(metric) {
			filtered = append(filtered, metric)
		}
	}

	return filtered
}

// FMetric is a metric that
type FMetric struct {
	name   string
	metric float64
}

// NewFMetric creates a new float64 metric
func NewFMetric(name string) *FMetric {
	return &FMetric{
		name: name,
	}
}

// Incr adds one to the metric
func (m *FMetric) Incr() {
	m.metric++
}

// Decr subtrcts one from the metric
func (m *FMetric) Decr() {
	m.metric--
}

// Add adds an abitrary value to the underlying metric
func (m *FMetric) Add(add float64) {
	m.metric += add
}

// MarshalJSON implements the marshaler interface
func (m *FMetric) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}{
		m.name,
		m.metric,
	})
}

// string implements the stringer interface
func (m *FMetric) String() string {
	return fmt.Sprintf("[%s] %.4f", m.name, m.metric)
}

// AddFMetrics adds metric values together and returns the result
func AddFMetrics(metrics ...FMetric) float64 {
	var total float64
	for _, m := range metrics {
		total += m.metric
	}
	return total
}

// FilterFMetrics takes a list of metrics and returns only those that satisfy the filter metric
func FilterFMetrics(metrics []FMetric, filter func(FMetric) bool) []FMetric {
	if filter == nil {
		return metrics
	}

	var filtered []FMetric
	for _, metric := range metrics {
		if filter(metric) {
			filtered = append(filtered, metric)
		}
	}

	return filtered
}

// GetName returns the name of the metric
func (m *FMetric) GetName() string {
	return m.name
}

// GetValue returns the value of the metric
func (m *FMetric) GetValue() float64 {
	return m.metric
}
