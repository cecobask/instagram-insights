package followdata

import (
	"github.com/cecobask/instagram-insights/pkg/instagram"
	"github.com/cecobask/instagram-insights/pkg/instagram/followdata"
	"github.com/spf13/cobra"
)

const CommandNameFollowing = "following"

func NewFollowingCommand() *cobra.Command {
	return &cobra.Command{
		Use:   CommandNameFollowing,
		Short: "Retrieve a list of users who you follow",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := instagram.NewOptions(instagram.OutputTable)
			return followdata.NewHandler().Following(opts)
		},
	}
}
