package root

import (
	"fmt"
	"os"

	"github.com/cecobask/instagram-insights/cmd/followdata"
	"github.com/cecobask/instagram-insights/cmd/information"
	"github.com/spf13/cobra"
)

const CommandNameInstagram = "instagram"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     fmt.Sprintf("%s [command]", CommandNameInstagram),
		Aliases: []string{"ig"},
		Short:   "Instagram Insights CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.AddCommand(
		information.NewRootCommand(),
		followdata.NewRootCommand(),
	)
	cmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})
	return cmd
}
