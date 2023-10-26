package runner

import (
	"fmt"
	"strings"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func collectors(collector v1alpha1.Collector) []string {
	var commands []string
	if collector.Events != nil {
		commands = append(commands, eventsCollector(*collector.Events))
	}
	if collector.PodLogs != nil {
		commands = append(commands, podLogsCollector(*collector.PodLogs))
	}
	return commands
}

func eventsCollector(collector v1alpha1.EventsCollector) string {
	var b strings.Builder
	b.WriteString("kubectl get events")
	if collector.Name != "" {
		fmt.Fprintf(&b, " %s", collector.Name)
	}
	if collector.Selector != "" {
		fmt.Fprintf(&b, " -l %s", collector.Selector)
	}
	ns := collector.Namespace
	// TODO
	if collector.Namespace == "" {
		ns = "$NAMESPACE"
	}
	fmt.Fprintf(&b, " -n %s", ns)
	return b.String()
}

func podLogsCollector(collector v1alpha1.PodLogsCollector) string {
	var b strings.Builder
	b.WriteString("kubectl logs --prefix")
	if collector.Name != "" {
		fmt.Fprintf(&b, " %s", collector.Name)
	}
	if collector.Selector != "" {
		fmt.Fprintf(&b, " -l %s", collector.Selector)
	}
	ns := collector.Namespace
	// TODO
	if collector.Namespace == "" {
		ns = "$NAMESPACE"
	}
	fmt.Fprintf(&b, " -n %s", ns)
	if len(collector.Container) > 0 {
		fmt.Fprintf(&b, " -c %s", collector.Container)
	} else {
		b.WriteString(" --all-containers")
	}
	if collector.Tail != nil {
		fmt.Fprintf(&b, " --tail=%d", *collector.Tail)
	}
	return b.String()
}
