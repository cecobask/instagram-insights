package instagram

import (
	"reflect"
	"testing"

	"github.com/spf13/pflag"
)

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
		Order  string
		Output string
		SortBy string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "succeeds to validate all options",
			fields: fields{
				Order:  OrderAsc,
				Output: OutputTable,
				SortBy: FieldTimestamp,
			},
			wantErr: false,
		},
		{
			name: "fails to validate order",
			fields: fields{
				Order: "invalid",
			},
			wantErr: true,
		},
		{
			name: "fails to validate output",
			fields: fields{
				Order:  OrderAsc,
				Output: "invalid",
			},
			wantErr: true,
		},
		{
			name: "fails to validate sort by",
			fields: fields{
				Order:  OrderAsc,
				Output: OutputTable,
				SortBy: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Options{
				Order:  tt.fields.Order,
				Output: tt.fields.Output,
				SortBy: tt.fields.SortBy,
			}
			if err := o.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewOptions(t *testing.T) {
	type args struct {
		flags *pflag.FlagSet
	}
	tests := []struct {
		name    string
		args    args
		want    *Options
		wantErr bool
	}{
		{
			name: "succeeds to create options",
			args: args{
				flags: func() *pflag.FlagSet {
					flags := pflag.NewFlagSet("", pflag.ExitOnError)
					flags.String(FlagOrder, OrderAsc, "")
					flags.String(FlagOutput, OutputTable, "")
					flags.String(FlagSortBy, FieldTimestamp, "")
					return flags
				}(),
			},
			want: &Options{
				Order:  OrderAsc,
				Output: OutputTable,
				SortBy: FieldTimestamp,
			},
			wantErr: false,
		},
		{
			name: "fails to find flag order",
			args: args{
				flags: pflag.NewFlagSet("", pflag.ExitOnError),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fails to find flag output",
			args: args{
				flags: func() *pflag.FlagSet {
					flags := pflag.NewFlagSet("", pflag.ExitOnError)
					flags.String(FlagOrder, OrderAsc, "")
					return flags
				}(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fails to find flag sort by",
			args: args{
				flags: func() *pflag.FlagSet {
					flags := pflag.NewFlagSet("", pflag.ExitOnError)
					flags.String(FlagOrder, OrderAsc, "")
					flags.String(FlagOutput, OutputTable, "")
					return flags
				}(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOptions(tt.args.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOptions() got = %v, want %v", got, tt.want)
			}
		})
	}
}
