package followdata

import (
	"github.com/cecobask/instagram-insights/pkg/instagram"
	"github.com/cecobask/instagram-insights/pkg/instagram/followdata"
	"github.com/spf13/cobra"
)

const CommandNameFollowing = "following"

func NewFollowingCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandNameFollowing,
		Short: "Retrieve a list of users who you follow",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := instagram.NewOptions(cmd.Flags())
			if err != nil {
				return err
			}
			if err = opts.Validate(); err != nil {
				return err
			}
			following, err := followdata.NewHandler().Following(opts)
			if err != nil {
				return err
			}
			cmd.Print(*following)
			return nil
		},
		DisableAutoGenTag: true,
	}
	addCommonFlags(cmd)
	return cmd
}
