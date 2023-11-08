package followdata

import (
	"github.com/cecobask/instagram-insights/pkg/instagram"
	"github.com/cecobask/instagram-insights/pkg/instagram/followdata"
	"github.com/spf13/cobra"
)

const CommandNameFollowers = "followers"

func NewFollowersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandNameFollowers,
		Short: "Retrieve a list of users who follow you",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := instagram.NewOptions(cmd.Flags())
			if err != nil {
				return err
			}
			if err = opts.Validate(); err != nil {
				return err
			}
			followers, err := followdata.NewHandler().Followers(opts)
			if err != nil {
				return err
			}
			cmd.Print(*followers)
			return nil
		},
		DisableAutoGenTag: true,
	}
	addCommonFlags(cmd)
	return cmd
}
