package mutate

func clean(in any) map[string]any {
	data, ok := in.(map[any]any)
	if !ok {
		return nil
	}
	out := map[string]any{}
	for k, v := range data {
		if c := clean(v); c != nil {
			out[k.(string)] = c
		} else {
			out[k.(string)] = v
		}
	}
	return out
}
