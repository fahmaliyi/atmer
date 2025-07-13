package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fahmaliyi/atmer/internal/service"
	"github.com/fahmaliyi/atmer/internal/utils"
	"github.com/spf13/cobra"
)

var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "Interactive menu to view, search, add, edit, or delete ATMs",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		for {
			machines, err := utils.LoadMachines(excelpath)
			if err != nil {
				fmt.Println("❌ Failed to load ATM list:", err)
				os.Exit(1)
			}

			fmt.Println("\n🔧 ATM Manager")
			fmt.Println("1. View all ATMs")
			fmt.Println("2. Search ATM by name")
			fmt.Println("3. Add new ATM")
			fmt.Println("4. Edit ATM IP")
			fmt.Println("5. Delete ATM")
			fmt.Println("6. Exit")
			fmt.Print("👉 Select an option (1-6): ")

			choice := readLine(reader)
			fmt.Println()

			switch choice {
			case "1":
				fmt.Printf("📄 Total ATMs: %d\n", len(machines))
				for i, m := range machines {
					fmt.Printf("%d. %s (%s)\n", i+1, m.Name, m.IP)
				}

			case "2":
				fmt.Print("🔍 Enter part of the ATM name: ")
				query := strings.ToLower(readLine(reader))
				var matches []service.Machine
				for _, m := range machines {
					if strings.Contains(strings.ToLower(m.Name), query) {
						matches = append(matches, m)
					}
				}
				if len(matches) == 0 {
					fmt.Println("❌ No matching ATMs found")
					break
				}
				fmt.Println("🔎 Matching ATMs:")
				for i, m := range matches {
					fmt.Printf("%d. %s (%s)\n", i+1, m.Name, m.IP)
				}

			case "3":
				fmt.Print("➕ Enter ATM name: ")
				name := readLine(reader)
				for _, m := range machines {
					if strings.EqualFold(m.Name, name) {
						fmt.Println("❌ ATM already exists")
						goto Menu
					}
				}
				var ip string
				for {
					fmt.Print("🔌 Enter ATM IP address: ")
					ip = readLine(reader)
					if isValidIP(ip) {
						break
					}
					fmt.Println("❌ Invalid IP format. Try again.")
				}
				machines = append(machines, service.Machine{Name: name, IP: ip})
				err := utils.SaveMachines(machines, excelpath)
				if err != nil {
					fmt.Println("❌ Failed to add ATM:", err)
				} else {
					fmt.Println("✅ ATM added")
				}

			case "4":
				fmt.Print("✏️ Enter part of ATM name to edit: ")
				query := strings.ToLower(readLine(reader))
				var matches []int
				for i, m := range machines {
					if strings.Contains(strings.ToLower(m.Name), query) {
						matches = append(matches, i)
					}
				}
				if len(matches) == 0 {
					fmt.Println("❌ No matching ATMs found")
					break
				}
				fmt.Println("🔎 Matching ATMs:")
				for idx, mi := range matches {
					m := machines[mi]
					fmt.Printf("%d. %s (%s)\n", idx+1, m.Name, m.IP)
				}

				var selected int
				for {
					fmt.Print("➡️ Select number to edit (or 'q' to cancel): ")
					input := readLine(reader)
					if input == "q" {
						fmt.Println("🚫 Cancelled")
						goto Menu
					}
					selectedIdx, err := strconv.Atoi(input)
					if err == nil && selectedIdx > 0 && selectedIdx <= len(matches) {
						selected = matches[selectedIdx-1]
						break
					}
					fmt.Println("❌ Invalid selection.")
				}

				fmt.Printf("✏️ Current IP: %s\n", machines[selected].IP)
				var newIP string
				for {
					fmt.Print("🔄 Enter new IP address: ")
					newIP = readLine(reader)
					if isValidIP(newIP) {
						break
					}
					fmt.Println("❌ Invalid IP format.")
				}
				machines[selected].IP = newIP
				err := utils.SaveMachines(machines, excelpath)
				if err != nil {
					fmt.Println("❌ Failed to save changes:", err)
				} else {
					fmt.Println("✅ ATM updated")
				}

			case "5":
				fmt.Print("🗑️ Enter part of ATM name to delete: ")
				query := strings.ToLower(readLine(reader))
				var matches []int
				for i, m := range machines {
					if strings.Contains(strings.ToLower(m.Name), query) {
						matches = append(matches, i)
					}
				}
				if len(matches) == 0 {
					fmt.Println("❌ No matching ATMs found")
					break
				}
				fmt.Println("🔎 Matching ATMs:")
				for idx, mi := range matches {
					m := machines[mi]
					fmt.Printf("%d. %s (%s)\n", idx+1, m.Name, m.IP)
				}

				var selected int
				for {
					fmt.Print("➡️ Select number to delete (or 'q' to cancel): ")
					input := readLine(reader)
					if input == "q" {
						fmt.Println("🚫 Cancelled")
						goto Menu
					}
					selectedIdx, err := strconv.Atoi(input)
					if err == nil && selectedIdx > 0 && selectedIdx <= len(matches) {
						selected = matches[selectedIdx-1]
						break
					}
					fmt.Println("❌ Invalid selection.")
				}

				target := machines[selected]
				fmt.Printf("❓ Confirm delete '%s' (y/N): ", target.Name)
				confirm := strings.ToLower(readLine(reader))
				if confirm != "y" && confirm != "yes" {
					fmt.Println("🚫 Cancelled")
					break
				}

				newList := append(machines[:selected], machines[selected+1:]...)
				err := utils.SaveMachines(newList, excelpath)
				if err != nil {
					fmt.Println("❌ Failed to delete ATM:", err)
				} else {
					fmt.Println("✅ ATM deleted")
				}

			case "6":
				fmt.Println("👋 Exiting ATM Manager.")
				return

			default:
				fmt.Println("❌ Invalid option")
			}

		Menu:
			fmt.Print("\n🔁 Press Enter to return to the menu...")
			readLine(reader)
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(manageCmd)
	manageCmd.Flags().StringVarP(&excelpath, "path", "p", "atms.xlsx", "Path to Excel file")
}

func readLine(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func isValidIP(ip string) bool {
	regex := `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}` +
		`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	match, _ := regexp.MatchString(regex, ip)
	return match
}
