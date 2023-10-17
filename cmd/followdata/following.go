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
			o, err := cmd.Flags().GetString(instagram.FlagOutput)
			if err != nil {
				return err
			}
			output, err := instagram.NewOutput(o)
			if err != nil {
				return err
			}
			opts := instagram.NewOptions(*output)
			following, err := followdata.NewHandler().Following(opts)
			if err != nil {
				return err
			}
			cmd.Print(*following)
			return nil
		},
		DisableAutoGenTag: true,
	}
}
