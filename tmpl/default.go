package tmpl

const Default = `
{{ define "changelog" }}
{{ if .PrevTag }}
# [{{ version }}]({{ url }}/compare/{{ .PrevTag }}..{{ .Tag }}) ({{ .Date }})
{{ else }}
# {{ version }} ({{ .Date }})
{{ end }}
{{ if .Breaking }}### Breaking Changes
{{ range .Breaking }}
{{ template "commit" . }}{{ end }}
{{ end }}
{{ if .Features }}### Features
{{ range .Features }}
{{ template "commit" . }}{{ end }}
{{ end }}
{{ if .BugFixes }}### Bug Fixes
{{ range .BugFixes }}
{{ template "commit" . }}{{ end }}{{ end }}
{{ end }}

{{ define "commit" }} * {{ .Subject }} ([{{ .ShortHash }}]({{ url }}/commit/{{ .ShortHash }})){{ end }}
`
