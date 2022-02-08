//go:build !windows
// +build !windows

package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"syscall"
)

func killGroup() {
	gid, err := syscall.Getpgid(os.Getpid())
	cobra.CheckErr(err)

	err = syscall.Kill(-gid, syscall.SIGINT)
	cobra.CheckErr(err)
}
