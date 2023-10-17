package instagram

import (
	"fmt"
)

type Options struct {
	Output string
}

func NewOptions(output string) Options {
	return Options{
		Output: output,
	}
}

func (o Options) Validate() error {
	return validateOutput(o.Output)
}

func validateOutput(value string) error {
	switch value {
	case OutputNone, OutputJson, OutputTable, OutputYaml:
		return nil
	default:
		return fmt.Errorf("invalid output format: %s", value)
	}
}
