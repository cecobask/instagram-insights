package instagram

import "testing"

func Test_validateOutputOption(t *testing.T) {
	type args struct {
		output string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "table output",
			args: args{
				output: OutputTable,
			},
			wantErr: false,
		},
		{
			name: "none output",
			args: args{
				output: OutputNone,
			},
			wantErr: false,
		},
		{
			name: "invalid output",
			args: args{
				output: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateOutputOption(tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("validateOutputOption() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
