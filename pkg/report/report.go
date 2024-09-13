package report

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/model"
)

func Save(report *model.Report, format v1alpha2.ReportFormatType, path, name string) error {
	if filepath.Ext(name) == "" {
		name += "." + strings.ToLower(string(format))
	}
	filePath := name
	if path != "" {
		filePath = filepath.Join(path, name)
	}
	switch format {
	case v1alpha2.XMLFormat:
		return saveJUnit(report, filePath)
	case v1alpha2.JSONFormat:
		return saveJson(report, filePath)
	default:
		return fmt.Errorf("unknown report format: %s", format)
	}
}
