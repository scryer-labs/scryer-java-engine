//go:build !darwin && !linux

package worker

import "os/exec"

func configureWorkerGroup(*exec.Cmd) {}
