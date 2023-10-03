package information

import (
	"fmt"
	"github.com/cecobask/instagram-insights/pkg/instagram"
	"github.com/spf13/cobra"
)

func NewDownloadCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "download <url>",
		Short: "Download instagram information locally",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("must provide exactly one archive download url")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := instagram.FetchInstagramInformation(args[0])
			if err != nil {
				return err
			}
			fmt.Println("downloaded instagram information locally")
			return nil
		},
	}
}
