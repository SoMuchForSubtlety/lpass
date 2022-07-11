/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/SoMuchForSubtlety/lpass/pkg/store"
	"github.com/SoMuchForSubtlety/lpass/pkg/util"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "l"},
	Short:   "List secrets",
	Run: func(cmd *cobra.Command, args []string) {
		entries, err := store.Load()
		util.HandleErr(err)

		if len(entries) == 0 {
			fmt.Println("no entries found, try the refresh command")
			return
		}

		for _, e := range entries {
			fmt.Println(e.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
