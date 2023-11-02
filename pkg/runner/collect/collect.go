package collect

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func Commands(collector *v1alpha1.Collect) ([]*v1alpha1.Command, error) {
	if collector == nil {
		return nil, nil
	}
	var commands []*v1alpha1.Command
	{
		cmd, err := podLogs(collector.PodLogs)
		if err != nil {
			return nil, err
		}
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	{
		cmd, err := events(collector.Events)
		if err != nil {
			return nil, err
		}
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands, nil
}
