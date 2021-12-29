package main

import (
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/frail/relmonorep/tmpl"

	"github.com/coreos/go-semver/semver"
)

type ChangeLog struct {
	Tag      string
	PrevTag  string
	Date     string
	Breaking []Commit
	Features []Commit
	BugFixes []Commit
	Others   []Commit
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

func generateChangelog(prefix, output string) error {
	git := Git{prefix}
	url, host, err := git.RepoInfo()
	if err != nil {
		return err
	}

	tags, err := git.ListReleaseTags()
	if err != nil {
		return fmt.Errorf("can't list release tags: %v", err)
	}

	if len(tags) == 0 {
		log.Fatal("There is no release tags to work on create a release tag first !")
	}

	lastReleaseTag := tags[0]
	c := ChangeLog{
		Tag:     "HEAD",
		PrevTag: lastReleaseTag,
		Date:    time.Now().Format("2006-01-02"),
	}
	if err := git.FillCommitsSince(&c); err != nil {
		return fmt.Errorf("can't read commit logs: %v", err)
	}

	lastReleaseVersion := lastReleaseTag
	if prefix != "" {
		lastReleaseVersion = lastReleaseTag[len(prefix)+1:]
	}

	version, err := semver.NewVersion(lastReleaseVersion)
	if err != nil {
		return err
	}

	switch {
	case len(c.Breaking) > 0:
		version.BumpMajor()
	case len(c.Features) > 0:
		version.BumpMinor()
	default:
		version.BumpPatch()
	}
	
	if prefix != "" {
		c.Tag = fmt.Sprintf("%s/%s", prefix, version)
	} else {
		c.Tag = version.String()
	}

	writer, err := newFile(output)
	if err != nil {
		return err
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

	var currentTemplate string
	switch host {
	case "github.com":
		currentTemplate = tmpl.Github
	case "bitbucket.org":
		currentTemplate = tmpl.BitBucket
	default:
		currentTemplate = tmpl.Default
	}

	t := template.Must(template.New("changelog").Funcs(functions).Parse(currentTemplate))
	if err := t.Execute(writer, c); err != nil {
		return err
	}

	return nil
}
