package template

const DefaultTemplate = true

func Get(values ...*bool) bool {
	for _, value := range values {
		if value != nil {
			return *value
		}
	}
	return DefaultTemplate
}
