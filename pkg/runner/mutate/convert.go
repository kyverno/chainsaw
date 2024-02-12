package mutate

func convert(in any) map[string]any {
	data, ok := in.(map[any]any)
	if !ok {
		return nil
	}
	out := map[string]any{}
	for k, v := range data {
		if c := convert(v); c != nil {
			out[k.(string)] = c
		} else {
			out[k.(string)] = v
		}
	}
	return out
}
