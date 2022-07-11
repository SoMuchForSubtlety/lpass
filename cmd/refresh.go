/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/SoMuchForSubtlety/lpass/pkg/api"
	"github.com/SoMuchForSubtlety/lpass/pkg/store"
	"github.com/SoMuchForSubtlety/lpass/pkg/util"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "refresh secrets from",
	Run: func(cmd *cobra.Command, args []string) {
		username, pw, err := store.LoadAPICredentials()
		if err != nil && errors.Is(err, store.ErrKeyNotFound) {
			username, pw = promptForCreds()
			err = store.StoreAPICrendentials(username, pw)
			if err != nil {
				logrus.Warnf("failed to persist LastPass credentials: %v", err)
			}
		}
		// TODO: only prompt if OTP is required
		otp := promptForOTP()
		entries, err := api.Load(context.Background(), username, pw, otp)
		if err != nil {
			if errors.Is(err, api.ErrBadCreds) {
				fmt.Println(err)
				err = store.DeleteAPICredentials()
				util.HandleErr(err)
			} else {
				util.HandleErr(err)
			}
		} else {
			err = store.Store(entries)
			util.HandleErr(err)
			fmt.Println("secrets refreshed")
		}
	},
}

func promptForCreds() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("please enter your LastPass username: ")
	username, err := reader.ReadString('\n')
	util.HandleErr(err)
	username = strings.TrimRight(username, "\n")
	fmt.Print("please enter your LastPass password: ")
	pass, err := reader.ReadString('\n')
	pass = strings.TrimRight(pass, "\n")
	util.HandleErr(err)
	return username, pass
}

func promptForOTP() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("please enter your LastPass OTP: ")
	otp, err := reader.ReadString('\n')
	otp = strings.TrimRight(otp, "\n")
	util.HandleErr(err)
	return otp
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}
