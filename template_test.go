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
				opts: []ErrOption{WithCode(1), FmtNoStack(true), FmtNoID(true)},
			},
			`{"code":1}`,
		},
		{
			"no stack or id template",
			fields{
				opts: []ErrOption{FmtNoStack(true), FmtNoID(true)},
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
			e := tr.New(tt.args.opts...)

			if got := e.Error(); got != tt.want {
				t.Errorf("Template.New() = %v, want %v", got, tt.want)
			}
		})
	}
}
