package cmd

import (
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
)

func Run(cmdline string) (string, error) {
	cmdArgs, err := shellwords.Parse(cmdline)

	if err != nil {
		return "", err
	}

	var cmd *exec.Cmd

	if len(cmdArgs) > 1 {
		cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
	} else {
		cmd = exec.Command(cmdArgs[0])
	}

	buf, err := cmd.CombinedOutput()
	out := strings.TrimSpace(string(buf))

	return out, err
}
