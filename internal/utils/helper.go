package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/fahmaliyi/atmer/internal/service"
	"github.com/xuri/excelize/v2"
)

func LoadMachines(path string) ([]service.Machine, error) {
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

func WriteResults(results []service.PingResult, output string) error {
	grouped := map[string][]string{
		"Offline":  {},
		"OnlyADSL": {},
	}

	for _, r := range results {
		if r.Status == "Offline" || r.Status == "OnlyADSL" {
			grouped[r.Status] = append(grouped[r.Status], r.Name)
		}
	}

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, status := range []string{"Offline", "OnlyADSL"} {
		names := grouped[status]
		if len(names) == 0 {
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
