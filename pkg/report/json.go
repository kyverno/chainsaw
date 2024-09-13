package report

import (
	"encoding/json"
	"os"

	"github.com/kyverno/chainsaw/pkg/model"
)

func saveJson(report *model.Report, file string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0o600)
}
