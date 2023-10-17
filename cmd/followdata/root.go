package followdata

import (
	"fmt"

	"github.com/cecobask/instagram-insights/pkg/instagram"
	"github.com/spf13/cobra"
)

const CommandNameFollowdata = "followdata"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     fmt.Sprintf("%s [command]", CommandNameFollowdata),
		Aliases: []string{"fd"},
		Short:   "Instagram follow data operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		DisableAutoGenTag: true,
	}
	cmd.AddCommand(NewFollowersCommand())
	cmd.AddCommand(NewFollowingCommand())
	cmd.AddCommand(NewUnfollowersCommand())
	for _, childCmd := range cmd.Commands() {
		childCmd.Flags().StringP(instagram.FlagOutput, "o", instagram.OutputTable, `output format ("json", "table", "yaml")`)
	}
	return cmd
}
