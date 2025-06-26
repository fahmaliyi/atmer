package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fahmaliyi/atmer/internal/service"
	xlsx "github.com/tealeg/xlsx/v3"
	"github.com/xuri/excelize/v2"
)

func LoadMachines(path string) ([]service.Machine, error) {
	path = filepath.Clean(path)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	file, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	rows, err := file.GetRows(file.GetSheetName(0))
	if err != nil {
		return nil, err
	}

	var machines []service.Machine
	for _, row := range rows {
		if len(row) >= 2 {
			name := strings.TrimSpace(row[0])
			ip := strings.TrimSpace(row[1])
			if name != "" && ip != "" {
				machines = append(machines, service.Machine{Name: name, IP: ip})
			}
		}
	}
	return machines, nil
}

func WriteResults(results []service.PingResult, output, format string, excludeOffline bool) error {
	switch format {
	case "json":
		return writeJSON(results, output, excludeOffline)
	case "csv":
		return writeCSV(results, output, excludeOffline)
	case "xlsx":
		return writeXLSX(results, output, excludeOffline)
	default:
		return writeTXT(results, output, excludeOffline)
	}
}

func writeJSON(results []service.PingResult, output string, excludeOffline bool) error {
	filtered := filterResults(results, excludeOffline)
	data, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(output, data, 0644)
}

func writeXLSX(results []service.PingResult, output string, excludeOffline bool) error {
	filtered := filterResults(results, excludeOffline)
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Results")
	if err != nil {
		return err
	}

	row := sheet.AddRow()
	row.AddCell().Value = "Name"
	row.AddCell().Value = "IP"
	row.AddCell().Value = "Status"

	for _, r := range filtered {
		row := sheet.AddRow()
		row.AddCell().Value = r.Name
		row.AddCell().Value = r.IP
		row.AddCell().Value = r.Status
	}

	return file.Save(output)
}

func writeCSV(results []service.PingResult, output string, excludeOffline bool) error {
	filtered := filterResults(results, excludeOffline)
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{"Name", "IP", "Status"})
	for _, r := range filtered {
		w.Write([]string{r.Name, r.IP, r.Status})
	}

	return nil
}

func writeTXT(results []service.PingResult, output string, excludeOffline bool) error {
	grouped := map[string][]string{
		"OnlyADSL": {},
	}
	// Include Offline only if excludeOffline is false
	if !excludeOffline {
		grouped["Offline"] = []string{}
	}

	for _, r := range results {
		if r.Status == "OnlyADSL" {
			grouped["OnlyADSL"] = append(grouped["OnlyADSL"], r.Name)
		} else if r.Status == "Offline" && !excludeOffline {
			grouped["Offline"] = append(grouped["Offline"], r.Name)
		}
	}

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	// Iterate statuses depending on presence in grouped
	for _, status := range []string{"Offline", "OnlyADSL"} {
		names, exists := grouped[status]
		if !exists || len(names) == 0 {
			continue
		}

		// Comma-separated
		f.WriteString(fmt.Sprintf("%s:\n", status))
		f.WriteString(strings.Join(names, ", ") + "\n\n")

		// Newline-separated
		f.WriteString(fmt.Sprintf("%s:\n", status))
		for _, name := range names {
			f.WriteString(name + "\n")
		}
		f.WriteString("\n")
	}

	return nil
}

func filterResults(results []service.PingResult, excludeOffline bool) []service.PingResult {
	if !excludeOffline {
		return results
	}
	filtered := make([]service.PingResult, 0, len(results))
	for _, r := range results {
		if r.Status != "Offline" {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
