package v1alpha2

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	conversion "k8s.io/apimachinery/pkg/conversion"
)

func Convert_v1alpha2_ConfigurationSpec_To_v1alpha1_ConfigurationSpec(in *ConfigurationSpec, out *v1alpha1.ConfigurationSpec, _ conversion.Scope) error {
	out.Catch = in.Catch
	out.Clusters = in.Clusters
	out.Timeouts = in.Timeouts
	if in := in.Cleanup; in != nil {
		out.SkipDelete = in.SkipDelete
		out.DelayBeforeCleanup = in.DelayBeforeCleanup
	}
	out.ExcludeTestRegex = in.Discovery.ExcludeTestRegex
	out.IncludeTestRegex = in.Discovery.IncludeTestRegex
	out.TestFile = in.Discovery.TestFile
	out.FullName = in.Discovery.FullName
	if in := in.Execution; in != nil {
		out.FailFast = in.FailFast
		out.Parallel = in.Parallel
		out.RepeatCount = in.RepeatCount
		out.ForceTerminationGracePeriod = in.ForceTerminationGracePeriod
	}
	if in := in.Namespace; in != nil {
		out.Namespace = in.Name
		out.NamespaceTemplate = in.Template
	}
	if in := in.Report; in != nil {
		out.ReportFormat = v1alpha1.ReportFormatType(in.Format)
		out.ReportPath = in.Path
		out.ReportName = in.Name
	}
	if in := in.Templating; in != nil {
		out.Template = in.Enabled
	}
	return nil
}

func Convert_v1alpha1_ConfigurationSpec_To_v1alpha2_ConfigurationSpec(in *v1alpha1.ConfigurationSpec, out *ConfigurationSpec, _ conversion.Scope) error {
	out.Catch = in.Catch
	out.Clusters = in.Clusters
	out.Timeouts = in.Timeouts
	out.Cleanup = &Cleanup{
		SkipDelete:         in.SkipDelete,
		DelayBeforeCleanup: in.DelayBeforeCleanup,
	}
	out.Discovery = Discovery{
		ExcludeTestRegex: in.ExcludeTestRegex,
		IncludeTestRegex: in.IncludeTestRegex,
		TestFile:         in.TestFile,
		FullName:         in.FullName,
	}
	out.Execution = &Execution{
		FailFast:                    in.FailFast,
		Parallel:                    in.Parallel,
		RepeatCount:                 in.RepeatCount,
		ForceTerminationGracePeriod: in.ForceTerminationGracePeriod,
	}
	out.Namespace = &Namespace{
		Name:     in.Namespace,
		Template: in.NamespaceTemplate,
	}
	out.Report = &Report{
		Format: ReportFormatType(in.ReportFormat),
		Path:   in.ReportPath,
		Name:   in.ReportName,
	}
	out.Templating = &Templating{
		Enabled: in.Template,
	}
	return nil
}
