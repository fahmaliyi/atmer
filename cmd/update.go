package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fahmaliyi/atmer/internal/service"
	"github.com/spf13/cobra"
)

var (
	updateFile  string
	updateKey   string
	updateVal   string
	updateMatch string
)

var updateServiceCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a service record field by service number",
	Run: func(cmd *cobra.Command, args []string) {
		content, err := os.ReadFile(updateFile)
		if err != nil {
			fmt.Printf("❌ Failed to read JSON file: %s\n", err)
			return
		}

		var records []service.ServiceRecord
		if err := json.Unmarshal(content, &records); err != nil {
			fmt.Printf("❌ Failed to parse JSON: %s\n", err)
			return
		}

		updated := false
		for i, r := range records {
			if strings.EqualFold(strings.TrimSpace(r.LANIP), strings.TrimSpace(updateMatch)) {
				switch strings.ToLower(updateKey) {
				case "location":
					records[i].Location = updateVal
				case "wanip":
					records[i].WANIP = updateVal
				case "lanip":
					records[i].LANIP = updateVal
				case "connectiontype":
					records[i].ConnectionType = updateVal
				case "bandwidth":
					records[i].Bandwidth = updateVal
				case "linetype":
					records[i].LineType = updateVal
				case "servicenumber":
					records[i].ServiceNumber = updateVal
				case "accountnumber":
					records[i].AccountNumber = updateVal
				default:
					fmt.Printf("❌ Unknown field: %s\n", updateKey)
					return
				}
			}

			updated = true
		}

		if !updated {
			fmt.Printf("❌ No record found with WANIP '%s'.\n", updateMatch)
			return
		}

		newContent, err := json.MarshalIndent(records, "", "  ")
		if err != nil {
			fmt.Printf("❌ Failed to marshal updated JSON: %s\n", err)
			return
		}

		if err := os.WriteFile(updateFile, newContent, 0644); err != nil {
			fmt.Printf("❌ Failed to write JSON file: %s\n", err)
			return
		}

		fmt.Printf("✅ Record with WANIP %s updated successfully.\n", updateMatch)
	},
}

func init() {
	// update flags
	updateServiceCmd.Flags().StringVarP(&updateFile, "file", "f", "services.json", "Path to JSON file")
	updateServiceCmd.Flags().StringVarP(&updateMatch, "match", "m", "", "Unique field value to identify the record (e.g. WANIP)")
	updateServiceCmd.Flags().StringVarP(&updateKey, "key", "k", "", "Field to update (location, wanip, lanip, etc.)")
	updateServiceCmd.Flags().StringVarP(&updateVal, "value", "v", "", "New value for the field")

	rootCmd.AddCommand(updateServiceCmd)
}
