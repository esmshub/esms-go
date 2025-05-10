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
	title := figure.NewFigure("ESMSgo!", "smslant", true)
	fmt.Println(strings.Repeat("_", 50))
	title.Print()
	fmt.Println(strings.Repeat("-", 50))
}

func init() {
	oldCmd := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		printBanner()
		// Then show the default help output
		oldCmd(cmd, args)
	})
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "esmsgo",
	Long: `ESMS Go! is a unified command-line interface for the ESMS fantasy football engine that 
wraps the core match engine and supporting utilities that were traditionally separate tools.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Configure(flagVerbose)
		if cmd.Name() != "help" {
			printBanner()
			fmt.Println("Command:", cmd.Short)
			fmt.Println(strings.Repeat("-", 50))
		}
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config-file", "c", "", "config file")
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "verbose output")

	// rootCmd.MarkPersistentFlagRequired("config-file")
}
