package information

import (
	"github.com/cecobask/instagram-insights/pkg/instagram/information"
	"github.com/spf13/cobra"
)

const CommandNameCleanup = "cleanup"

func NewCleanupCommand() *cobra.Command {
	return &cobra.Command{
		Use:   CommandNameCleanup,
		Short: "Cleanup local instagram information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return information.NewHandler().Cleanup()
		},
		DisableAutoGenTag: true,
	}
}
