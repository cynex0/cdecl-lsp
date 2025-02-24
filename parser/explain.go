package parser

import "os/exec"

func Explain(line string) (string, error) {
	cmd := exec.Command("cdecl", "explain", line)
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(stdout), nil
}
