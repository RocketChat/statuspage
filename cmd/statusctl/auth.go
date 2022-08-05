package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/RocketChat/statuscentral/cmd/statusctl/common"
)

var baseURL = ""
var loginToken = ""

var loginCmd = &cobra.Command{
	Use:     "login",
	Short:   "Login",
	Example: "statusctl login",
	Run: func(c *cobra.Command, args []string) {
		if err := common.Login(baseURL, loginToken); err != nil {
			panic(err)
		}

		log.Println("Login Complete!")
	},
}

var logoutCmd = &cobra.Command{
	Use:     "logout",
	Short:   "logout",
	Example: "statusctl logout",
	Run: func(c *cobra.Command, args []string) {
		if err := common.Logout(); err != nil {
			panic(err)
		}
	},
}

func init() {
	loginCmd.Flags().StringVarP(&baseURL, "host", "", "", "status page host")
	loginCmd.Flags().StringVarP(&loginToken, "token", "t", "", "auth token")
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
}
