/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/SoMuchForSubtlety/lpass/pkg/store"
	"github.com/SoMuchForSubtlety/lpass/pkg/ui"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get SECRET_NAME",
	Short: "Search for secrets matching the provided query",
	RunE:  get,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		entries, err := store.Load()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveNoFileComp
		}

		var matching []string
		for _, e := range entries {
			if strings.HasPrefix(strings.ToLower(e.Name), strings.ToLower(toComplete)) {
				matching = append(matching, e.Name)
			}
		}

		return matching, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
	},
}

var regex *bool

func get(cmd *cobra.Command, args []string) error {
	entries, err := store.Load()
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return errors.New("no entries found, try the refresh command")

	}

	if len(args) == 0 {
		return errors.New("please provide a query")
	}
	query := args[0]

	matchFn := func(entry string) bool {
		query = strings.ToLower(query)
		return strings.Contains(strings.ToLower(entry), query)
	}
	if *regex {
		re, err := regexp.Compile(query)
		if err != nil {
			return fmt.Errorf("invalid regexp: %w", err)
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
		return nil
	} else if len(matches) == 0 {
		fmt.Println("no matching entries found")
		return nil
	}

	entry := selectEntry(matches)
	if entry == nil {
		return nil
	}
	ui.Render(*entry)
	return nil
}

func selectEntry(entries []store.Entry) *store.Entry {
	return ui.Select(entries)
}

func init() {
	rootCmd.AddCommand(getCmd)
	regex = getCmd.Flags().BoolP("regex", "r", false, "Interpret the query as a regular expression")
}
