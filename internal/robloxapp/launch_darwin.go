//go:build darwin

package robloxapp

import (
	"errors"
	"os"
	"path/filepath"
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
