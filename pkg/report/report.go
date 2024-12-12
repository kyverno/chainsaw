package report

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/model"
)

func getFile(path, name, extension string) string {
	if filepath.Ext(name) == "" {
		name += "." + strings.ToLower(extension)
	}
	filePath := name
	if path != "" {
		filePath = filepath.Join(path, name)
	}
	return filePath
}

func Save(report *model.Report, format v1alpha2.ReportFormatType, path, name string) error {
	switch format {
	case v1alpha2.XMLFormat, v1alpha2.JUnitTestFormat:
		return saveJUnitTest(report, getFile(path, name, "xml"))
	case v1alpha2.JUnitStepFormat:
		return saveJUnitStep(report, getFile(path, name, "xml"))
	case v1alpha2.JUnitOperationFormat:
		return saveJUnitOperation(report, getFile(path, name, "xml"))
	case v1alpha2.JSONFormat:
		return saveJson(report, getFile(path, name, "json"))
	default:
		return fmt.Errorf("unknown report format: %s", format)
	}
}
