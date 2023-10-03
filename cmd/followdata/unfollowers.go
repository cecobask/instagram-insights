package followdata

import (
	"github.com/cecobask/instagram-insights/pkg/instagram"
	"github.com/spf13/cobra"
)

func NewUnfollowersCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "unfollowers",
		Short: "Find out which instagram users are not following back",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
}

func run() error {
	followData := instagram.NewFollowData()
	if err := followData.ExtractAllData(); err != nil {
		return err
	}
	followData.FindUnfollowers()
	return nil
}
