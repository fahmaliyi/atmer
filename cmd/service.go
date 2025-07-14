package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fahmaliyi/atmer/internal/service"
	"github.com/fahmaliyi/atmer/internal/utils"
	"github.com/fatih/color"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/spf13/cobra"
)

var searchTerm string
var serviceFile string

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Search service details by partial match on any field",
	Run: func(cmd *cobra.Command, args []string) {

		green := color.New(color.FgGreen, color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		if searchTerm == "" {
			fmt.Println("❌ Please provide a search term using -s")
			return
		}

		content, err := os.ReadFile(serviceFile)
		if err != nil {
			fmt.Printf("❌ Failed to read JSON file: %s\n", err)
			return
		}

		var records []service.ServiceRecord
		if err := json.Unmarshal(content, &records); err != nil {
			fmt.Printf("❌ Failed to parse JSON: %s\n", err)
			return
		}

		query := strings.TrimSpace(strings.ToLower(searchTerm))
		var matches []service.ServiceRecord

		for _, r := range records {
			bandwidthStr := fmt.Sprintf("%v", r.Bandwidth)

			fields := []string{
				r.Location,
				r.WANIP,
				r.LANIP,
				r.ConnectionType,
				bandwidthStr,
				r.LineType,
				fmt.Sprintf("%v", r.ServiceNumber),
				fmt.Sprintf("%v", r.AccountNumber),
			}

			for _, field := range fields {
				if fuzzy.Match(query, strings.ToLower(field)) {
					matches = append(matches, r)
					break
				}
			}
		}

		if len(matches) == 0 {
			fmt.Println("❌ No matches found.")
			return
		}

		fmt.Printf("🔍 Found %d matching service(s):\n\n", len(matches))
		for _, r := range matches {
			bandwidthStr := fmt.Sprintf("%v", r.Bandwidth)

			// Header line: location and WAN IP
			fmt.Printf("%s %s\n", green("- "+r.Location), green("("+r.WANIP+")"))

			// Details lines with colored labels and values
			fmt.Printf("  %s %s | %s %s | %s %s | %s %s\n",
				cyan("LAN:"), yellow(r.LANIP),
				cyan("Conn:"), yellow(r.ConnectionType),
				cyan("BW:"), yellow(bandwidthStr),
				cyan("Line:"), yellow(r.LineType),
			)
			fmt.Printf("  %s %s | %s %s\n\n",
				cyan("Service #:"), yellow(utils.ToString(r.ServiceNumber)),
				cyan("Account #:"), yellow(utils.ToString(r.AccountNumber)),
			)

		}
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.Flags().StringVarP(&searchTerm, "search", "s", "", "Search term (partial match)")
	serviceCmd.Flags().StringVarP(&serviceFile, "file", "f", "services.json", "Path to JSON file")
}
