/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	figure "github.com/common-nighthawk/go-figure"
	"github.com/esmshub/esms-go/cmd/cli/logger"
	"github.com/spf13/cobra"
)

var flagVerbose bool

func printBanner() {
	tagline := "Electronic Soccer Management Simulator v3.4"
	title := figure.NewFigure("ESMS", "smslant", true)
	fmt.Println(title.String())
	fmt.Println(tagline)
	fmt.Println(strings.Repeat("~", len(tagline)))
	fmt.Println()
}

func init() {
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		printBanner()
		// Then show the default help output
		cmd.Usage()
	})
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "esmscli",
	Short: "A modern CLI that brings together the ESMS match engine and its legacy utilities into a single, streamlined tool.",
	Long: `A unified command-line interface for the ESMS fantasy soccer simulator, consolidating the core match engine and a variety of supporting utilities that were traditionally separate tools.

Originally developed in the late 1990s, ESMS was the foundation for hundreds of fantasy soccer leagues, run via play-by-email. Over time, various tools emerged to support simulation, stat tracking, team management, and configuration — often shared informally and maintained independently.

This CLI brings those components under one roof, modernizing the experience while staying true to the spirit of the original system. Whether you're running a simulation, configuring your league, or generating outputs, esmscli provides a consistent and extensible interface designed for both nostalgia and usability.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Configure(flagVerbose)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config-file", "c", "", "config file (default is $HOME/.cli.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "verbose output")

	// rootCmd.MarkPersistentFlagRequired("config-file")
}
