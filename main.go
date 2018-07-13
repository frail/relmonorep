package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var url = flag.String("repo", "", "github url of a repo")
var re = regexp.MustCompile(`^(?P<hash>[0-9a-fA-F]{40}) ((?P<type>fix|feat|perf)(?P<project>\(\w+\))?: )?(?P<subject>.+)$`)

func main() {
	flag.Parse()
	functions := template.FuncMap{
		"url": func() string {
			return *url
		},
	}
	t := template.Must(template.New("changelog").Funcs(functions).Parse(ChangeLogTemplate))
	// c := ChangeLog{
	// 	Tag:     "release/1.1.0",
	// 	PrevTag: "release/1.0.0",
	// 	Date:    "2018-07-13",
	// 	BugFixes: []Commit{
	// 		{"0782afb26f57f58a03f97404174d76de", "Define functions on structs, not pointers", "", ""},
	// 	},
	// 	Features: []Commit{
	// 		{"ebedf17e60fefa28eea6b381b4a05718", "Add url function", "", ""},
	// 		{"de4ff5e78ace8a8bf89172c33891b957", "Initial implementation", "", ""},
	// 	},
	// }

	git := Git{"release/playbook/"}
	tags, _ := git.ListReleaseTags()
	c := ChangeLog{
		Tag:     "HEAD",
		PrevTag: tags[0],
		Date:    time.Now().Format("2006-01-02"),
	}
	git.FillCommitsSince(&c)
	t.Execute(os.Stdout, c)
}

type Git struct {
	Prefix string
}

func (g *Git) FillCommitsSince(cl *ChangeLog) error {
	lines, err := g.exec("log", "--format=%H %s", fmt.Sprintf("%s..%s", cl.PrevTag, cl.Tag))
	if err != nil {
		return err
	}
	for _, line := range lines {
		c := Commit{}
		res := re.FindStringSubmatch(line)
		for i, name := range re.SubexpNames() {
			switch name {
			case "hash":
				c.Hash = res[i]
			case "type":
				c.Type = res[i]
			case "project":
				// TODO remove parens
				c.Project = res[i]
			case "subject":
				c.Subject = res[i]
			}
		}
		switch c.Type {
		case "fix", "perf":
			cl.BugFixes = append(cl.BugFixes, c)
		case "feat":
			cl.Features = append(cl.Features, c)
		default:
			cl.Others = append(cl.Others, c)
		}
	}
	return nil
}

func (g *Git) ListReleaseTags() ([]string, error) {
	return g.exec("tag", "-l", g.Prefix+"*", "--sort=-creatordate")
}

func (g *Git) exec(args ...string) ([]string, error) {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(out))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

const ChangeLogTemplate = `
{{ define "changelog" }}
{{ if .PrevTag }}
# [{{ .Version }}]({{ url }}/compare/{{ .PrevTag }}..{{ .Tag }}) ({{ .Date }})
{{ else }}
# {{ .Version }} ({{ .Date }})
{{ end }}

### Bug Fixes
{{ range .BugFixes }}
{{ template "commit" . }}{{ end }}

### Features
{{ range .Features }}
{{ template "commit" . }}{{ end }}
{{ end }}

{{ define "commit" }} * {{ .Subject }} ([{{ .ShortHash }}]({{ url }}/commit/{{ .ShortHash }})) {{ end }}
`

type ChangeLog struct {
	Tag      string
	PrevTag  string
	Date     string
	Features []Commit
	BugFixes []Commit
	Others   []Commit
}

func (c ChangeLog) Version() string {
	return c.Tag[strings.LastIndexByte(c.Tag, '/')+1:]
}

type Commit struct {
	Hash    string
	Subject string
	Type    string
	Project string
}

func (c Commit) ShortHash() string {
	return c.Hash[:7]
}
