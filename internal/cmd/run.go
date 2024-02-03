package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
)

func Run(cmdline string) (string, error) {
	envs, args, err := shellwords.ParseWithEnvs(cmdline)

	if err != nil {
		return "", err
	}

	var cmd *exec.Cmd

	if len(args) > 1 {
		cmd = exec.Command(args[0], args[1:]...)
	} else {
		cmd = exec.Command(args[0])
	}

	if len(envs) > 0 {
		cmd.Env = append(os.Environ(), envs...)
	}

	buf, err := cmd.CombinedOutput()
	out := strings.TrimSpace(string(buf)) + "\n"

	return out, err
}
