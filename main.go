package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/coreos/go-semver/semver"
)

var prefix = flag.String("prefix", "", "release branch prefix")

func main() {
	flag.Parse()
	git := Git{*prefix}
	url, err := git.RepoURL()
	if err != nil {
		panic(err)
	}

	tags, err := git.ListReleaseTags()
	if err != nil {
		panic(err)
	}

	lastRelease := tags[0]
	c := ChangeLog{
		Tag:     "HEAD",
		PrevTag: lastRelease,
		Date:    time.Now().Format("2006-01-02"),
	}
	if err := git.FillCommitsSince(&c); err != nil {
		panic(err)
	}

	version, err := semver.NewVersion(lastRelease[len(*prefix)+1:])
	if err != nil {
		panic(err)
	}

	switch {
	case len(c.Breaking) > 0:
		version.BumpMajor()
	case len(c.Features) > 0:
		version.BumpMinor()
	default:
		version.BumpPatch()
	}
	c.Tag = fmt.Sprintf("%s/%s", *prefix, version)

	functions := template.FuncMap{
		"url": func() string {
			return url
		},
		"version": func() string {
			return version.String()
		},
	}
	t := template.Must(template.New("changelog").Funcs(functions).Parse(ChangeLogTemplate))
	if err := t.Execute(os.Stdout, c); err != nil {
		panic(err)
	}
}
