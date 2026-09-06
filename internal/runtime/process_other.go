//go:build !windows

package runtime

import "os/exec"

func configureCommand(cmd *exec.Cmd) {}
