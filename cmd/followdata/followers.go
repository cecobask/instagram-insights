package followdata

import (
	"github.com/cecobask/instagram-insights/pkg/instagram"
	"github.com/cecobask/instagram-insights/pkg/instagram/followdata"
	"github.com/spf13/cobra"
)

const CommandNameFollowers = "followers"

func NewFollowersCommand() *cobra.Command {
	return &cobra.Command{
		Use:   CommandNameFollowers,
		Short: "Retrieve a list of users who follow you",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := instagram.NewOptions(instagram.OutputTable)
			return followdata.NewHandler().Followers(opts)
		},
	}
}
