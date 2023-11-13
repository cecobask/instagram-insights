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
	cmd.AddCommand(
		NewFollowersCommand(),
		NewFollowingCommand(),
		NewUnfollowersCommand(),
	)
	return cmd
}

func addCommonFlags(cmd *cobra.Command) {
	cmd.Flags().Int(instagram.FlagLimit, instagram.Unlimited, `maximum results to display, leave empty for unlimited`)
	cmd.Flags().String(instagram.FlagOrder, instagram.OrderDesc, `order direction ("asc", "desc")`)
	cmd.Flags().String(instagram.FlagOutput, instagram.OutputTable, `output format ("json", "table", "yaml")`)
	cmd.Flags().String(instagram.FlagSortBy, instagram.FieldTimestamp, `sort by field ("timestamp", "username")`)
}
