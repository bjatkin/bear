package bear

import "testing"

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
