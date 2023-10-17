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
			output, err := cmd.Flags().GetString(instagram.FlagOutput)
			if err != nil {
				return err
			}
			opts := instagram.NewOptions(output)
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
}
