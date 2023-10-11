package followdata

import (
	"github.com/spf13/cobra"
)

const CommandNameFollowdata = "followdata"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandNameFollowdata,
		Short: "Instagram follow data operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(NewUnfollowersCommand())
	return cmd
}
