package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/RocketChat/statuscentral/cmd/statusctl/common"
)

var loginCmd = &cobra.Command{
	Use:     "login",
	Short:   "OAuth Login",
	Example: "statusctl login",
	Run: func(c *cobra.Command, args []string) {
		if err := common.Login(); err != nil {
			panic(err)
		}

		log.Println("Login Complete!")
	},
}

var logoutCmd = &cobra.Command{
	Use:     "logout",
	Short:   "OAuth logout",
	Example: "statusctl logout",
	Run: func(c *cobra.Command, args []string) {
		if err := common.Logout(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
}
