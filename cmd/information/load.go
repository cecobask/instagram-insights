package information

import (
	"fmt"
	"strings"

	"github.com/cecobask/instagram-insights/pkg/instagram/information"
	"github.com/spf13/cobra"
)

const CommandNameLoad = "load"

func NewLoadCommand() *cobra.Command {
	return &cobra.Command{
		Use: CommandNameLoad + " <source>",
		Example: func() string {
			examples := []string{
				"instagram information load https://drive.google.com/file/d/xyz",
				"instagram information load file:///home/username/Desktop/instagram_data.zip",
			}
			return strings.Join(examples, "\n")
		}(),
		Short: "Load Instagram information",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("must provide exactly one location")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return information.NewHandler().Load(args[0])
		},
		DisableAutoGenTag: true,
	}
}
