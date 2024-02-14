package template

const DefaultTemplate = false

func Get(values ...*bool) bool {
	for _, value := range values {
		if value != nil {
			return *value
		}
	}
	return DefaultTemplate
}
