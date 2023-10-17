package instagram

import (
	"reflect"
	"testing"
)

func TestNewOutput(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    *Output
		wantErr bool
	}{
		{
			name: "succeeds to create output",
			args: args{
				value: "json",
			},
			want: func() *Output {
				o := OutputJson
				return &o
			}(),
			wantErr: false,
		},
		{
			name: "fails to create output",
			args: args{
				value: "invalid",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOutput(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOutput() got = %v, want %v", got, tt.want)
			}
		})
	}
}
