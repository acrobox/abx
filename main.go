package main

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"

	"acrobox.io/abx/cli"
)

func main() {
	config := &cli.Config{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	config.Home = os.Getenv("ACROBOX_HOME")
	if config.Home == "" {
		home, err := os.UserConfigDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
			return
		}
		config.Home = filepath.Join(home, cli.AppName)
	}
	err := cli.Run(config)
	if err != nil {
		serr, ok := err.(*ssh.ExitError)
		if ok {
			status := serr.ExitStatus()
			os.Exit(status)
			return
		}
		os.Exit(1)
	}
}
