package v1alpha2

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	conversion "k8s.io/apimachinery/pkg/conversion"
)

func Convert_v1alpha2_ConfigurationSpec_To_v1alpha1_ConfigurationSpec(in *ConfigurationSpec, out *v1alpha1.ConfigurationSpec, _ conversion.Scope) error {
	out.SkipDelete = in.Cleanup.SkipDelete
	out.DelayBeforeCleanup = in.Cleanup.DelayBeforeCleanup
	out.Clusters = in.Clusters
	out.DeletionPropagationPolicy = in.Deletion.Propagation
	out.ExcludeTestRegex = in.Discovery.ExcludeTestRegex
	out.IncludeTestRegex = in.Discovery.IncludeTestRegex
	out.TestFile = in.Discovery.TestFile
	out.FullName = in.Discovery.FullName
	// TODO
	// if in := in.Error; in != nil {
	// 	out.Catch = in.Catch
	// }
	out.FailFast = in.Execution.FailFast
	out.Parallel = in.Execution.Parallel
	out.RepeatCount = in.Execution.RepeatCount
	out.ForceTerminationGracePeriod = in.Execution.ForceTerminationGracePeriod
	out.Namespace = in.Namespace.Name
	out.NamespaceTemplate = in.Namespace.Template
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
	out.Cleanup = CleanupOptions{
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
	// TODO
	// out.Error = &ErrorOptions{
	// 	Catch: in.Catch,
	// }
	out.Execution = ExecutionOptions{
		FailFast:                    in.FailFast,
		Parallel:                    in.Parallel,
		RepeatCount:                 in.RepeatCount,
		ForceTerminationGracePeriod: in.ForceTerminationGracePeriod,
	}
	out.Namespace = NamespaceOptions{
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
