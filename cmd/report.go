package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fahmaliyi/atmer/internal/service"
	"github.com/fahmaliyi/atmer/internal/utils"
	"github.com/spf13/cobra"
)

var (
	excelpath string
	noOffline bool
	noOnline  bool
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate ATM connectivity report from Excel file",
	Run: func(cmd *cobra.Command, args []string) {
		machines, err := utils.LoadMachines(excelpath)
		if err != nil {
			fmt.Println("❌ Failed to load:", err)
			os.Exit(1)
		}

		var results []service.PingResult

		for _, m := range machines {
			status := "Offline"
			if service.Ping(m.IP) {
				status = "Online"
			} else if service.Ping(service.GetModemIP(m.IP)) {
				status = "OnlyADSL"
			}

			if noOnline && status == "Online" {
				continue // Skip offline machines if --no-offline is set
			}

			fmt.Printf("- %s (%s) → %s\n", m.Name, m.IP, status)

			results = append(results, service.PingResult{
				Name:   m.Name,
				IP:     m.IP,
				Status: status,
			})
		}

		output, _ := cmd.Flags().GetString("output")

		ext := strings.ToLower(filepath.Ext(output))
		var format string

		switch ext {
		case ".json":
			format = "json"
		case ".csv":
			format = "csv"
		case ".xlsx", ".xls":
			format = "xlsx"
		default:
			format = "txt"
		}

		if format == "" {
			fmt.Printf("⚠️ Unknown output format for extension '%s'. Defaulting to txt.\n", ext)
			format = "txt"
		}

		err = utils.WriteResults(results, output, format, noOffline)
		if err != nil {
			fmt.Println("❌ Error writing results:", err)
		} else {
			fmt.Println("✅ Results written to", output)
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringP("output", "o", "ping_results.txt", "Output file")
	reportCmd.Flags().StringP("format", "f", "txt", "Output format: txt, json, csv, xlsx")
	reportCmd.Flags().StringVarP(&excelpath, "path", "p", "atms.xlsx", "Path to Excel file")
	reportCmd.Flags().BoolVar(&noOffline, "no-offline", false, "Exclude offline ATMs from report")
	reportCmd.Flags().BoolVar(&noOnline, "no-online", false, "Exclude online ATMs from report")
}
