/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/SoMuchForSubtlety/lpass/pkg/store"
	"github.com/SoMuchForSubtlety/lpass/pkg/ui"
	"github.com/SoMuchForSubtlety/lpass/pkg/util"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Search for secrets matching the provided query",
	Run:   get,
}

var regex *bool

func get(cmd *cobra.Command, args []string) {
	entries, err := store.Load()
	util.HandleErr(err)
	if len(entries) == 0 {
		fmt.Println("no entries found, try the refresh command")
		return
	}

	if len(args) == 0 {
		fmt.Println("please provide a query")
		os.Exit(1)
	}
	query := args[0]

	matchFn := func(entry string) bool {
		query = strings.ToLower(query)
		return strings.Contains(strings.ToLower(entry), query)
	}
	if *regex {
		re, err := regexp.Compile(query)
		if err != nil {
			util.HandleErr(fmt.Errorf("invalid regexp: %w", err))
		}
		matchFn = func(entry string) bool {
			return re.MatchString(entry)
		}
	}

	var matches []store.Entry

	for _, entry := range entries {
		if matchFn(entry.Name) || matchFn(entry.Notes) || matchFn(entry.ID) {
			matches = append(matches, entry)
		}
	}
	if len(matches) == 1 {
		ui.Render(matches[0])
		return
	} else if len(matches) == 0 {
		fmt.Println("no matching entries found")
		return
	}

	entry := selectEntry(matches)
	if entry == nil {
		return
	}
	ui.Render(*entry)
}

func selectEntry(entries []store.Entry) *store.Entry {
	return ui.Select(entries)
}

func init() {
	rootCmd.AddCommand(getCmd)
	regex = getCmd.Flags().BoolP("regex", "r", false, "Interpret the query as a regular expression")
}
