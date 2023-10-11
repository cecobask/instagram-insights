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

func validateOutputOption(output string) error {
	switch output {
	case OutputTable:
		return nil
	case OutputNone:
		return nil
	default:
		return fmt.Errorf("invalid output format: %s", output)
	}
}
