{{- define "comment" -}}
  {{- $comment := "" -}}
  {{- range . -}}
    {{- if . -}}
      {{- if not (eq (index . 0) '+') -}}
        {{- if $comment -}}
          {{- $comment = print $comment " " . -}}
        {{- else -}}
          {{- $comment = . -}}
        {{- end -}}
      {{- end -}}
    {{- end -}}
  {{- end -}}
  {{- if $comment -}}
    <p>{{ $comment }}</p>
  {{- else -}}
    *No description provided.*
  {{- end -}}
{{- end -}}

{{- define "typ" -}}
  {{- if .Link -}}
    [`{{ .DisplayName }}`]({{ .Link }})
  {{- else -}}
    `{{ .DisplayName }}`
  {{- end -}}
{{- end -}}

{{- define "members" }}
  {{- range .GetMembers }}
    {{- if not .Hidden }}
      {{- $name := .FieldName }}
      {{- $optional := .IsOptional }}
      {{- $type := .GetType }}
      {{- $inline := .IsInline }}
      {{- $comment := .GetComment }}
| `{{ $name }}` | {{ template "typ" $type }} | {{ if not $optional }}:white_check_mark:{{ end }} | {{ template "comment" .CommentLines }} |
    {{- end }}
  {{- end }}
{{- end }}
