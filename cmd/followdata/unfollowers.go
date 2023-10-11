package followdata

import (
	"github.com/cecobask/instagram-insights/pkg/instagram"
	"github.com/cecobask/instagram-insights/pkg/instagram/followdata"
	"github.com/spf13/cobra"
)

const CommandNameUnfollowers = "unfollowers"

func NewUnfollowersCommand() *cobra.Command {
	return &cobra.Command{
		Use:   CommandNameUnfollowers,
		Short: "Find out which instagram users are not following back",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := instagram.NewOptions(instagram.OutputTable)
			return followdata.NewHandler().Unfollowers(opts)
		},
	}
}
