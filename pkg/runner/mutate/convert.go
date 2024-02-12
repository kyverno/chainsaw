package mutate

func convertSlice(in any) []any {
	data, ok := in.([]any)
	if !ok {
		return nil
	}
	out := []any{}
	for _, v := range data {
		out = append(out, convert(v))
	}
	return out
}

func convertMap(in any) map[string]any {
	data, ok := in.(map[any]any)
	if !ok {
		return nil
	}
	out := map[string]any{}
	for k, v := range data {
		out[k.(string)] = convert(v)
	}
	return out
}

func convert(in any) any {
	if result := convertSlice(in); result != nil {
		return result
	}
	if result := convertMap(in); result != nil {
		return result
	}
	return in
}
