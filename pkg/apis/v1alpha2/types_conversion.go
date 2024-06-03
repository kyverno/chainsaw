package v1alpha2

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	conversion "k8s.io/apimachinery/pkg/conversion"
)

func Convert_v1alpha2_ConfigurationSpec_To_v1alpha1_ConfigurationSpec(in *ConfigurationSpec, out *v1alpha1.ConfigurationSpec, _ conversion.Scope) error {
	if in := in.Cleanup; in != nil {
		out.SkipDelete = in.SkipDelete
		out.DelayBeforeCleanup = in.DelayBeforeCleanup
	}
	out.Clusters = in.Clusters
	out.DeletionPropagationPolicy = in.Deletion.Propagation
	out.ExcludeTestRegex = in.Discovery.ExcludeTestRegex
	out.IncludeTestRegex = in.Discovery.IncludeTestRegex
	out.TestFile = in.Discovery.TestFile
	out.FullName = in.Discovery.FullName
	if in := in.Error; in != nil {
		out.Catch = in.Catch
	}
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
	out.Template = &in.Templating.Enabled
	out.Timeouts = in.Timeouts
	return nil
}

func Convert_v1alpha1_ConfigurationSpec_To_v1alpha2_ConfigurationSpec(in *v1alpha1.ConfigurationSpec, out *ConfigurationSpec, _ conversion.Scope) error {
	out.Cleanup = &CleanupOptions{
		SkipDelete:         in.SkipDelete,
		DelayBeforeCleanup: in.DelayBeforeCleanup,
	}
	out.Clusters = in.Clusters
	out.Deletion = DeletionOptions{
		Propagation: in.DeletionPropagationPolicy,
	}
	out.Discovery = DiscoveryOptions{
		ExcludeTestRegex: in.ExcludeTestRegex,
		IncludeTestRegex: in.IncludeTestRegex,
		TestFile:         in.TestFile,
		FullName:         in.FullName,
	}
	out.Error = &ErrorOptions{
		Catch: in.Catch,
	}
	out.Execution = &ExecutionOptions{
		FailFast:                    in.FailFast,
		Parallel:                    in.Parallel,
		RepeatCount:                 in.RepeatCount,
		ForceTerminationGracePeriod: in.ForceTerminationGracePeriod,
	}
	out.Namespace = &NamespaceOptions{
		Name:     in.Namespace,
		Template: in.NamespaceTemplate,
	}
	out.Report = &ReportOptions{
		Format: ReportFormatType(in.ReportFormat),
		Path:   in.ReportPath,
		Name:   in.ReportName,
	}
	if in.Template == nil {
		out.Templating = TemplatingOptions{
			Enabled: true,
		}
	} else {
		out.Templating = TemplatingOptions{
			Enabled: *in.Template,
		}
	}
	out.Timeouts = in.Timeouts
	return nil
}
