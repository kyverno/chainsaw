package cleanup

func Skip(config bool, test *bool, step *bool) bool {
	if step != nil {
		return *step
	}
	if test != nil {
		return *test
	}
	return config
}
