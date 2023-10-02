package information

import "github.com/spf13/cobra"

type Options struct {
	ArchiveDownloadURL string
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "information",
		Short: "Instagram information operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(NewDownloadCommand())
	cmd.AddCommand(NewCleanupCommand())
	return cmd
}
