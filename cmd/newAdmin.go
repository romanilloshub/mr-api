/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"orov.io/siempreAbierto/plugin/auth"
)

// newAdminCmd represents the newAdmin command
var newAdminCmd = &cobra.Command{
	Use:   "newAdmin",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		authUtil := auth.NewUtil(context.Background())
		err := authUtil.NewAdmin(userEmail)
		if err != nil {
			log.Printf("Unable to create user with email %s due to: %s", userEmail, err)
		}
		log.Printf("User %s create succesfully!", userEmail)
	},
}

var userEmail string

func init() {
	userCmd.AddCommand(newAdminCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newAdminCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	newAdminCmd.Flags().StringVarP(&userEmail, "email", "e", "", "Help message for toggle")
	newAdminCmd.MarkFlagRequired("email")
}
