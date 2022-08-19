package metrics

import (
	"reflect"
	"testing"
)

func TestMetric_Incr(t *testing.T) {
	m := NewMetric("Test Metric")
	if m.GetValue() != 0 {
		t.Fatalf("Incr() wanted initial value of 0 but got %d", m.GetValue())
	}

	m.Incr()
	if m.GetValue() != 1 {
		t.Fatalf("Incr() wanted initial value of 1 but got %d", m.GetValue())
	}

	m.Incr()
	if m.GetValue() != 2 {
		t.Fatalf("Incr() wanted initial value of 2 but got %d", m.GetValue())
	}
}

func TestMetric_Decr(t *testing.T) {
	m := NewMetric("Test Metric")
	if m.GetValue() != 0 {
		t.Fatalf("Incr() wanted initial value of 0 but got %d", m.GetValue())
	}

	m.Decr()
	if m.GetValue() != -1 {
		t.Fatalf("Incr() wanted initial value of -1 but got %d", m.GetValue())
	}

	m.Decr()
	if m.GetValue() != -2 {
		t.Fatalf("Incr() wanted initial value of -2 but got %d", m.GetValue())
	}
}

func TestMetric_Add(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		add int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			"add one",
			fields{
				"Test Metric",
			},
			args{
				add: 1,
			},
			1,
		},
		{
			"add 10",
			fields{
				"Test Metric",
			},
			args{
				add: 10,
			},
			10,
		},
		{
			"subtract 5",
			fields{
				"Test Metric",
			},
			args{
				add: -5,
			},
			-5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetric(tt.fields.name)
			m.Add(tt.args.add)

			if m.GetValue() != tt.want {
				t.Fatalf("Add() got %d, want %d", m.GetValue(), tt.want)
			}
		})
	}
}

func TestMetric_String(t *testing.T) {
	type fields struct {
		name   string
		metric int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"new",
			fields{
				name: "Test Metric",
			},
			"[Test Metric] 0",
		},
		{
			"updated metric",
			fields{
				name:   "Test Metric",
				metric: 50,
			},
			"[Test Metric] 50",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metric{
				name:   tt.fields.name,
				metric: tt.fields.metric,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("Metric.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddMetrics(t *testing.T) {
	type args struct {
		metrics []Metric
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"single metric",
			args{
				metrics: []Metric{{name: "A", metric: 10}},
			},
			10,
		},
		{
			"two metrics",
			args{
				metrics: []Metric{{name: "A", metric: 10}, {name: "B", metric: 5}},
			},
			15,
		},
		{
			"three metrics",
			args{
				metrics: []Metric{{name: "A", metric: 10}, {name: "B", metric: 5}, {name: "C", metric: 7}},
			},
			22,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddMetrics(tt.args.metrics...); got != tt.want {
				t.Errorf("AddMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterMetrics(t *testing.T) {
	type args struct {
		metrics []Metric
		filter  func(Metric) bool
	}
	tests := []struct {
		name string
		args args
		want []Metric
	}{
		{
			"nil filter function",
			args{
				metrics: []Metric{{name: "A"}},
				filter:  nil,
			},
			[]Metric{{name: "A"}},
		},
		{
			"no more metrics",
			args{
				metrics: []Metric{{name: "A"}, {name: "B"}},
				filter: func(Metric) bool {
					return false
				},
			},
			nil,
		},
		{
			"all metrics",
			args{
				metrics: []Metric{{name: "A"}, {name: "B"}},
				filter: func(Metric) bool {
					return true
				},
			},
			[]Metric{{name: "A"}, {name: "B"}},
		},
		{
			"some metrics",
			args{
				metrics: []Metric{{name: "A", metric: 5}, {name: "B", metric: -5}, {name: "C"}},
				filter: func(m Metric) bool {
					if m.GetName() == "A" {
						return true
					}
					if m.GetValue() < 0 {
						return true
					}

					return false
				},
			},
			[]Metric{{name: "A", metric: 5}, {name: "B", metric: -5}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterMetrics(tt.args.metrics, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetric_MarshalJSON(t *testing.T) {
	type fields struct {
		name   string
		metric float64
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"valid json",
			fields{
				name:   "Test Metric",
				metric: 12.44,
			},
			`{"name":"Test Metric","value":12.44}`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &FMetric{
				name:   tt.fields.name,
				metric: tt.fields.metric,
			}
			got, err := m.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Metric.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Metric.MarshalJSON() = %s, want %s", string(got), tt.want)
			}
		})
	}
}

// Test for FMetrics

func TestFMetric_Incr(t *testing.T) {
	m := NewFMetric("Test Metric")
	if m.GetValue() != 0 {
		t.Fatalf("Incr() wanted initial value of 0 but got %f", m.GetValue())
	}

	m.Incr()
	if m.GetValue() != 1 {
		t.Fatalf("Incr() wanted initial value of 1 but got %f", m.GetValue())
	}

	m.Incr()
	if m.GetValue() != 2 {
		t.Fatalf("Incr() wanted initial value of 2 but got %f", m.GetValue())
	}
}

func TestFMetric_Decr(t *testing.T) {
	m := NewFMetric("Test Metric")
	if m.GetValue() != 0 {
		t.Fatalf("Incr() wanted initial value of 0 but got %f", m.GetValue())
	}

	m.Decr()
	if m.GetValue() != -1 {
		t.Fatalf("Incr() wanted initial value of -1 but got %f", m.GetValue())
	}

	m.Decr()
	if m.GetValue() != -2 {
		t.Fatalf("Incr() wanted initial value of -2 but got %f", m.GetValue())
	}
}

func TestFMetric_Add(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		add float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			"add 1.1",
			fields{
				"Test Metric",
			},
			args{
				add: 1.1,
			},
			1.1,
		},
		{
			"add 10",
			fields{
				"Test Metric",
			},
			args{
				add: 10.55,
			},
			10.55,
		},
		{
			"subtract 5",
			fields{
				"Test Metric",
			},
			args{
				add: -5.678,
			},
			-5.678,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewFMetric(tt.fields.name)
			m.Add(tt.args.add)

			if m.GetValue() != tt.want {
				t.Fatalf("Add() got %f, want %f", m.GetValue(), tt.want)
			}
		})
	}
}

func TestFMetric_String(t *testing.T) {
	type fields struct {
		name   string
		metric float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"new",
			fields{
				name: "Test Metric",
			},
			"[Test Metric] 0.0000",
		},
		{
			"updated metric",
			fields{
				name:   "Test Metric",
				metric: 50.53839,
			},
			"[Test Metric] 50.5384",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &FMetric{
				name:   tt.fields.name,
				metric: tt.fields.metric,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("Metric.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddFMetrics(t *testing.T) {
	type args struct {
		metrics []FMetric
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"single metric",
			args{
				metrics: []FMetric{{name: "A", metric: 10.5}},
			},
			10.5,
		},
		{
			"two metrics",
			args{
				metrics: []FMetric{{name: "A", metric: 10.2}, {name: "B", metric: 5.678}},
			},
			15.878,
		},
		{
			"three metrics",
			args{
				metrics: []FMetric{{name: "A", metric: 10.2}, {name: "B", metric: 5.678}, {name: "C", metric: 7.01}},
			},
			22.888,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddFMetrics(tt.args.metrics...); (got - tt.want) > 0.0001 {
				t.Errorf("AddMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterFMetrics(t *testing.T) {
	type args struct {
		metrics []FMetric
		filter  func(FMetric) bool
	}
	tests := []struct {
		name string
		args args
		want []FMetric
	}{
		{
			"nil filter function",
			args{
				metrics: []FMetric{{name: "A"}},
				filter:  nil,
			},
			[]FMetric{{name: "A"}},
		},
		{
			"no more metrics",
			args{
				metrics: []FMetric{{name: "A"}, {name: "B"}},
				filter: func(FMetric) bool {
					return false
				},
			},
			nil,
		},
		{
			"all metrics",
			args{
				metrics: []FMetric{{name: "A"}, {name: "B"}},
				filter: func(FMetric) bool {
					return true
				},
			},
			[]FMetric{{name: "A"}, {name: "B"}},
		},
		{
			"some metrics",
			args{
				metrics: []FMetric{{name: "A", metric: 5.5}, {name: "B", metric: -5.5}, {name: "C"}},
				filter: func(m FMetric) bool {
					if m.GetName() == "A" {
						return true
					}
					if m.GetValue() < 0 {
						return true
					}

					return false
				},
			},
			[]FMetric{{name: "A", metric: 5.5}, {name: "B", metric: -5.5}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterFMetrics(tt.args.metrics, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFMetric_MarshalJSON(t *testing.T) {
	type fields struct {
		name   string
		metric float64
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"valid json",
			fields{
				name:   "Test Metric",
				metric: 10.5393,
			},
			`{"name":"Test Metric","value":10.5393}`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &FMetric{
				name:   tt.fields.name,
				metric: tt.fields.metric,
			}
			got, err := m.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Metric.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Metric.MarshalJSON() = %s, want %s", string(got), tt.want)
			}
		})
	}
}
