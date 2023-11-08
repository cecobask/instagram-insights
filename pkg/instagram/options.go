package instagram

import (
	"fmt"

	"github.com/spf13/pflag"
)

type Options struct {
	Order  string
	Output string
	SortBy string
}

func NewOptions(flags *pflag.FlagSet) (*Options, error) {
	order, err := flags.GetString(FlagOrder)
	if err != nil {
		return nil, err
	}
	output, err := flags.GetString(FlagOutput)
	if err != nil {
		return nil, err
	}
	sortBy, err := flags.GetString(FlagSortBy)
	if err != nil {
		return nil, err
	}
	return &Options{
		Order:  order,
		Output: output,
		SortBy: sortBy,
	}, nil
}

func NewEmptyOptions() *Options {
	return &Options{
		Output: OutputNone,
	}
}

func (o *Options) Validate() error {
	if err := validateOrder(o.Order); err != nil {
		return err
	}
	if err := validateOutput(o.Output); err != nil {
		return err
	}
	if err := validateSortBy(o.SortBy); err != nil {
		return err
	}
	return nil
}

func validateOrder(value string) error {
	switch value {
	case OrderAsc, OrderDesc:
		return nil
	default:
		return fmt.Errorf("invalid order direction: %s", value)
	}
}

func validateOutput(value string) error {
	switch value {
	case OutputNone, OutputJson, OutputTable, OutputYaml:
		return nil
	default:
		return fmt.Errorf("invalid output format: %s", value)
	}
}

func validateSortBy(value string) error {
	switch value {
	case FieldTimestamp, FieldUsername:
		return nil
	default:
		return fmt.Errorf("invalid sort by field: %s", value)
	}
}
