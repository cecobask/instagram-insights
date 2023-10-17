package instagram

import (
	"fmt"
)

type Options struct {
	Output Output
}

func NewOptions(o Output) Options {
	return Options{
		Output: o,
	}
}

type Output string

func NewOutput(value string) (*Output, error) {
	switch value {
	case string(OutputNone), string(OutputJson), string(OutputTable), string(OutputYaml):
		o := Output(value)
		return &o, nil
	default:
		return nil, fmt.Errorf("invalid output format: %s", value)
	}
}
