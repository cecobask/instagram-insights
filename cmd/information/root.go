package information

import "github.com/spf13/cobra"

const CommandNameInformation = "information"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandNameInformation,
		Short: "Instagram information operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(NewDownloadCommand())
	cmd.AddCommand(NewCleanupCommand())
	return cmd
}
