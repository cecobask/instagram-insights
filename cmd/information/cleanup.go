package information

import (
	"fmt"
	"github.com/cecobask/instagram-insights/pkg/instagram"
	"github.com/spf13/cobra"
)

func NewCleanupCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cleanup",
		Short: "Cleanup local instagram information",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := instagram.CleanupInstagramInformation()
			if err != nil {
				return err
			}
			fmt.Println("cleaned up local instagram information")
			return nil
		},
	}
}
