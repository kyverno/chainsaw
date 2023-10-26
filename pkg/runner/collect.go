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
	// TODO
	// if len(tc.Selector) > 0 {
	// 	fmt.Fprintf(&b, " -l %s", tc.Selector)
	// }
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
	// TODO
	// if len(tc.Selector) > 0 {
	// 	fmt.Fprintf(&b, " -l %s", tc.Selector)
	// }
	ns := collector.Namespace
	// TODO
	if collector.Namespace == "" {
		ns = "$NAMESPACE"
	}
	fmt.Fprintf(&b, " -n %s", ns)
	b.WriteString(" --all-containers")
	// TODO
	// if len(tc.Container) > 0 {
	// 	fmt.Fprintf(&b, " -c %s", tc.Container)
	// } else {
	// 	b.WriteString(" --all-containers")
	// }
	// if tc.Tail == 0 {
	// 	if len(tc.Selector) > 0 {
	// 		tc.Tail = 10
	// 	} else {
	// 		tc.Tail = -1
	// 	}
	// }
	// fmt.Fprintf(&b, " --tail=%d", tc.Tail)
	return b.String()
}
