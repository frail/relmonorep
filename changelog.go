package main

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
