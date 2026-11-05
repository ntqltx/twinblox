package main

import (
	"github.com/ntqltx/twinblox/internal/macosapp"
	"github.com/ntqltx/twinblox/internal/robloxapp"
	"github.com/ntqltx/twinblox/internal/urlhandler"
)

const BUNDLE_IDENTIFIER string = "com.ntqltx.twinblox"

func main() {
	if !urlhandler.Set(BUNDLE_IDENTIFIER, "roblox-player") || !urlhandler.Check(BUNDLE_IDENTIFIER, "roblox-player") {
		macosapp.Notification("can't make a handler for roblox binary links")
		return
	}

	macosapp.Run(func(url string) {
		if _, err := robloxapp.Launch(url); err != nil {
			macosapp.Notification("can't launch roblox: " + err.Error())
		}
	}, func() {
		urlhandler.Set(urlhandler.ROBLOX_BUNDLE_IDENTIFIER, "roblox-player")
	})
}
