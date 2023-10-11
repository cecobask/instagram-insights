package information

import (
	"fmt"

	"github.com/cecobask/instagram-insights/pkg/instagram/information"
	"github.com/spf13/cobra"
)

const CommandNameDownload = "download"

func NewDownloadCommand() *cobra.Command {
	return &cobra.Command{
		Use:   CommandNameDownload + " <url>",
		Short: "Download instagram information locally",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("must provide exactly one archive download url")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return information.NewHandler().Download(args[0])
		},
	}
}
