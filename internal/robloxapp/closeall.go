package robloxapp

import (
	"os"

	"github.com/ntqltx/twinblox/pkg/ps"
)

func CloseAll() {
	processes, err := ps.Processes()
	if err != nil {
		return
	}

	for _, process := range processes {
		pName := process.Executable()
		if pName == "RobloxPlayer" {
			p, err := os.FindProcess(process.Pid())
			if err == nil {
				p.Kill()
			}
		}
	}
}
