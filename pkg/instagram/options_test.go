package instagram

import "testing"

func Test_validateOutput(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid output",
			args: args{
				value: OutputNone,
			},
			wantErr: false,
		},
		{
			name: "invalid output",
			args: args{
				value: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateOutput(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("validateOutput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOptions_Validate(t *testing.T) {
	type fields struct {
		Output string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "succeeds to validate options",
			fields: fields{
				Output: OutputNone,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Options{
				Output: tt.fields.Output,
			}
			if err := o.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
