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
		Short: "Retrieve a list of users who are not following you back",
		RunE: func(cmd *cobra.Command, args []string) error {
			o, err := cmd.Flags().GetString(instagram.FlagOutput)
			if err != nil {
				return err
			}
			output, err := instagram.NewOutput(o)
			if err != nil {
				return err
			}
			opts := instagram.NewOptions(*output)
			unfollowers, err := followdata.NewHandler().Unfollowers(opts)
			if err != nil {
				return err
			}
			cmd.Print(*unfollowers)
			return nil
		},
		DisableAutoGenTag: true,
	}
}
