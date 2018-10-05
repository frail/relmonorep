package main

import (
	"flag"
	"fmt"
	"text/template"
	"time"

	"github.com/coreos/go-semver/semver"
)

var (
	prefix = flag.String("prefix", "", "release branch prefix")
	output = flag.String("output", "-", "output file to prepend changelog")
)

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

	lastReleaseTag := tags[0]
	c := ChangeLog{
		Tag:     "HEAD",
		PrevTag: lastReleaseTag,
		Date:    time.Now().Format("2006-01-02"),
	}
	if err := git.FillCommitsSince(&c); err != nil {
		panic(err)
	}

	lastReleaseVersion := lastReleaseTag
	if *prefix != "" {
		lastReleaseVersion = lastReleaseTag[len(*prefix)+1:]
	}

	version, err := semver.NewVersion(lastReleaseVersion)
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

	writer, err := newFile(*output)
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	functions := template.FuncMap{
		"url": func() string {
			return url
		},
		"version": func() string {
			return version.String()
		},
	}
	t := template.Must(template.New("changelog").Funcs(functions).Parse(ChangeLogTemplate))
	if err := t.Execute(writer, c); err != nil {
		panic(err)
	}
}
