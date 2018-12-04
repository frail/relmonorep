package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
)

var (
	commitRE = regexp.MustCompile(`^(?P<hash>[0-9a-fA-F]{40}) ((?P<type>fix|feat|perf|style|chore|refactor)(?P<project>\(\w+\))?: )?(?P<subject>.+)$`)
	sshRE    = regexp.MustCompile(`^git@(?P<host>[a-z.]+?):(?P<path>[0-9a-zA-Z./-]+?)\.git$`)
	versionRE = regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+$`)
)

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
		res := commitRE.FindStringSubmatch(line)
		for i, name := range commitRE.SubexpNames() {
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
		case "fix", "perf", "style":
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
	tags, err := g.exec("tag", "-l", g.Prefix+"*", "--sort=-creatordate")
	if err != nil {
		return nil, err
	}

	var relTags []string

	for _, t := range(tags) {
		if versionRE.MatchString(t) {
			relTags = append(relTags, t)
		}
	}
	return tags, err
}

func (g *Git) RepoURL() (string, error) {
	var url, host, path string
	remote, err := g.exec("remote", "get-url", "origin")
	if err != nil {
		return url, err
	}

	if len(remote) != 1 {
		return url, errors.New("git: more than 1 origin repo url")
	}

	remoteURL := remote[0]
	if sshRE.MatchString(remoteURL) {
		sshRes := sshRE.FindStringSubmatch(remoteURL)
		for i, name := range sshRE.SubexpNames() {
			switch name {
			case "host":
				host = sshRes[i]
			case "path":
				path = sshRes[i]
			}
		}
	}
	// TODO add https origin support

	url = fmt.Sprintf("https://%s/%s", host, path)
	return url, nil
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
