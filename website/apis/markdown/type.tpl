{{- define "type" }}
## `{{ .Name.Name }}`     {#{{ .Anchor }}}
{{- if eq .Kind "Alias" }}

(Alias of `{{ .Underlying }}`)
{{- end }}
{{- with .References }}

**Appears in:**
    {{ range . }}
        {{- if or .Referenced .IsExported }}
- [{{ .DisplayName }}]({{ .Link }})
        {{- end }}
    {{- end }}
{{- end }}
{{- if .GetComment }}

{{ .GetComment }}
{{- end }}
{{- if .GetMembers }}

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
    {{- /* . is a apiType */}}
    {{- if .IsExported }}
        {{- /* Add apiVersion and kind rows if deemed necessary */}}
| `apiVersion` | `string` | :white_check_mark: | | `{{- .APIGroup -}}` |
| `kind` | `string` | :white_check_mark: | | `{{- .Name.Name -}}` |
    {{- end }}
    {{- /* The actual list of members is in the following template */}}
    {{- template "members" . }}
    {{- end }}
{{ end }}
