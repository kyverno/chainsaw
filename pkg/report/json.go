package report

import (
	"encoding/json"
	"os"
)

func saveJson(report *Report, file string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0o600)
}
