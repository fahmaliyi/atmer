package cmd

import (
	"fmt"
	"os"

	"github.com/fahmcode/atmer/internal/service"
	"github.com/fahmcode/atmer/internal/utils"
	"github.com/spf13/cobra"
)

var (
	excelpath string
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

			fmt.Printf("- %s (%s) → %s\n", m.Name, m.IP, status)

			results = append(results, service.PingResult{
				Name:   m.Name,
				IP:     m.IP,
				Status: status,
			})
		}

		output, _ := cmd.Flags().GetString("output")
		err = utils.WriteResults(results, output)
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
	reportCmd.Flags().StringVarP(&excelpath, "path", "p", "atms.xlsx", "Path to Excel file")
}
