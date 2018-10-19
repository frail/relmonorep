package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

const releaseScript = `
set -ex

git checkout develop
git diff --exit-code >/dev/null
git pull
git diff --exit-code >/dev/null
git checkout master
git pull
git diff --exit-code >/dev/null
git merge develop


git tag "%s"
git push
git push --tags

git checkout develop
`

func release(prefix, version string, verbose bool) error {
	cmd := exec.Command("bash")
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
	}

	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	if prefix != "" {
		version = fmt.Sprintf("%s/%s", prefix, version)
	}
	script := fmt.Sprintf(releaseScript, version)
	if _, err := io.WriteString(in, script); err != nil {
		return err
	}

	if err = in.Close(); err != nil {
		return err
	}

	return cmd.Wait()
}
