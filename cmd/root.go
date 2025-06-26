package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "atmer",
	Short: "CLI tool to ping ATMs and generate connectivity reports from Excel data",
	Long: `Atmer is a command-line application designed to help network administrators and
support teams quickly assess the connectivity status of a list of ATMs by pinging their IP addresses.

Given an Excel file containing ATM names and their primary IP addresses, Atmer:
  - Pings each ATM's primary IP to check if it is online
  - Pings the secondary/modem IP (calculated automatically) to detect 'OnlyADSL' connectivity
  - Classifies ATMs as Online, OnlyADSL, or Offline
  - Generates a detailed report saved to a text file, grouping ATMs by status

This tool accelerates troubleshooting and network health monitoring for ATM fleets
with minimal setup, leveraging concurrency for fast performance.

Usage Examples:

  atmer report -p atms.xlsx -o ping_results.txt

Flags:

  -p, --path      Path to the Excel file with ATM data (default "atms.xlsx")
  -o, --output    Output file path for the generated report (default "ping_results.txt")

Atmer is built with Go and Cobra for reliable and efficient CLI experience.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
