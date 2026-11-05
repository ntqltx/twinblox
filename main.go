package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ntqltx/twinblox/internal/macosapp"
	"github.com/ntqltx/twinblox/internal/robloxapp"
	"github.com/ntqltx/twinblox/internal/urlhandler"
)

const BUNDLE_IDENTIFIER string = "com.ntqltx.twinblox"

// roblox apparently using both of these links, so we take over both
var schemes = []string{"roblox-player", "roblox"}

func main() {
	// twinblox opens roblox links instead of roblox application
	for _, s := range schemes {
		if !urlhandler.Set(BUNDLE_IDENTIFIER, s) || !urlhandler.Check(BUNDLE_IDENTIFIER, s) {
			macosapp.Notification("can't make a handler for roblox binary links")

			restoreHandlers()
			return
		}
	}

	// we're handling links restoring by checking all signals
	// if the instances gets closed/killed by user without quitting it normally through an app
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-sig
		restoreHandlers()
		os.Exit(0)
	}()

	macosapp.Run(func(url string) {
		if _, err := robloxapp.Launch(url); err != nil {
			macosapp.Notification("can't launch roblox: " + err.Error())
		}
	}, restoreHandlers)
}

// restoring all of the links back to Roblox.app
func restoreHandlers() {
	for _, s := range schemes {
		urlhandler.Set(urlhandler.ROBLOX_BUNDLE_IDENTIFIER, s)
	}
}
