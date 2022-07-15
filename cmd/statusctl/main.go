package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/RocketChat/statuscentral/buildInfo"
	"github.com/RocketChat/statuscentral/cmd/statusctl/incident"
	"github.com/RocketChat/statuscentral/cmd/statusctl/maintenance"
)

var rootCmd = &cobra.Command{
	Use:   "statusctl",
	Short: "StatusCentral cli",
}

func main() {
	rootCmd.Version = fmt.Sprintf("%s-alpha", buildInfo.GetVersion())
	rootCmd.AddCommand(incident.IncidentCmd)
	rootCmd.AddCommand(maintenance.MaintenanceCmd)
	rootCmd.Execute() //nolint:errcheck // Tech debt
}
