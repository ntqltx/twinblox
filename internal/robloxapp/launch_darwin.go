//go:build darwin

package robloxapp

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/ntqltx/twinblox/internal/syncbreaker"
)

const BINARY_PATH = "Contents/MacOS/RobloxPlayer"

var ErrNotInstalled = errors.New("Roblox.app not found in ~/Applications")

func AppPath() (string, error) {
	candidates := []string{"/Applications/Roblox.app"}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, "Applications", "Roblox.app"))
	}

	for _, c := range candidates {
		if _, err := os.Stat(filepath.Join(c, BINARY_PATH)); err == nil {
			return c, nil
		}
	}

	return "", ErrNotInstalled
}

func Launch(url string) (*os.Process, error) {
	app, err := AppPath()
	if err != nil {
		return nil, err
	}

	// we're skipping LaunchServices by executing the binary directly
	// so macos doesn't block any other roblox instances from opening
	cmd := exec.Command(filepath.Join(app, BINARY_PATH), "-protocolString", url)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	syncbreaker.Break() // we'll still delete semaphore just in case
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go cmd.Wait()
	return cmd.Process, nil
}
