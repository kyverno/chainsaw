package discovery

import (
	"os"
	"regexp"
	"slices"
)

var StepFileName = regexp.MustCompile(`^(\d\d)-(.*)\.(?:yaml|yml)$`)

type Step struct {
	AssertFiles []string
	ErrorFiles  []string
	OtherFiles  []string
}

func TryFindStepFiles(path string) (map[string]Step, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	} else {
		// collect and sort candidate files
		var stepFiles []string
		for _, file := range files {
			fileName := file.Name()
			if !file.IsDir() && StepFileName.MatchString(fileName) {
				stepFiles = append(stepFiles, fileName)
			}
		}
		if len(stepFiles) == 0 {
			return nil, nil
		}
		slices.Sort(stepFiles)
		steps := map[string]Step{}
		for _, file := range stepFiles {
			groups := StepFileName.FindStringSubmatch(file)
			s := steps[groups[1]]
			switch groups[2] {
			case "assert":
				s.AssertFiles = append(s.AssertFiles, file)
			case "errors":
				s.ErrorFiles = append(s.ErrorFiles, file)
			default:
				s.OtherFiles = append(s.OtherFiles, file)
			}
			steps[groups[1]] = s
		}
		return steps, nil
	}
}
