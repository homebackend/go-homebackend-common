package homecommon

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Execute(sudo bool, check bool, command []string) (int, string, string) {
	if sudo {
		command = append([]string{"sudo"}, command...)
	}

	cmd := exec.Command(command[0], command[1:]...)
	var cout bytes.Buffer
	var cerr bytes.Buffer
	cmd.Stdout = &cout
	cmd.Stderr = &cerr
	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			log.Printf("Exit Status for `%s`: %d (%s)", strings.Join(command, " "), exiterr.ExitCode(), cerr.String())
			if check {
				os.Exit(1)
			}
			return exiterr.ExitCode(), cout.String(), cerr.String()
		} else {
			log.Printf("cmd.Run: %v", err)
			return -1, "", "Other error"
		}
	}

	return 0, cout.String(), cerr.String()
}
