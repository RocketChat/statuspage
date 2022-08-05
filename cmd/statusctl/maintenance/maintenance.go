package maintenance

import (
	"fmt"

	"github.com/spf13/cobra"
)

var SubCommands []*cobra.Command

var outputFormat = "list"

var MaintenanceCmd = &cobra.Command{
	Use: "maintenance",
	Aliases: []string{
		"maint",
		"m",
	},
	Short:   "StatusCentral scheduled maintenance",
	Example: "statusctl maintenance [command]",
	Args: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("%v requires arguments", c.UseLine())
		}

		return nil
	},
}

func init() {
	getCmd.Flags().StringVarP(&outputFormat, "output", "o", "list", "output format")
	listCmd.Flags().BoolVarP(&latestOnly, "latest", "l", false, "Show latest only")

	SubCommands = append(SubCommands, listCmd, describeCmd, getCmd, createCmd, updateCmd)
	MaintenanceCmd.AddCommand(SubCommands...)
}
