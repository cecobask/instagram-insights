package information

import (
	"fmt"

	"github.com/spf13/cobra"
)

const CommandNameInformation = "information"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [command]", CommandNameInformation),
		Short: "Instagram information operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		DisableAutoGenTag: true,
	}
	cmd.AddCommand(NewDownloadCommand())
	cmd.AddCommand(NewCleanupCommand())
	return cmd
}
