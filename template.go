package main

const ChangeLogTemplate = `
{{ define "changelog" }}
{{ if .PrevTag }}
# [{{ version }}]({{ url }}/compare/{{ .PrevTag }}..{{ .Tag }}) ({{ .Date }})
{{ else }}
# {{ version }} ({{ .Date }})
{{ end }}

### Features
{{ range .Features }}
{{ template "commit" . }}{{ end }}

### Bug Fixes
{{ range .BugFixes }}
{{ template "commit" . }}{{ end }}
{{ end }}

{{ define "commit" }} * {{ .Subject }} ([{{ .ShortHash }}]({{ url }}/commit/{{ .ShortHash }})){{ end }}
`
