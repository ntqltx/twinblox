package main

import (
	"github.com/ntqltx/twinblox/internal/deeplink"
	"github.com/ntqltx/twinblox/internal/macosapp"
	"github.com/ntqltx/twinblox/internal/robloxapp"
	"github.com/ntqltx/twinblox/internal/urlhandler"
)

func main() {
	urlhandler.Set("com.insadem.multiroblox", "roblox-player")
	s := urlhandler.Check(urlhandler.ROBLOX_BUNDLE_IDENTIFIER, "roblox-player")
	if s {
		macosapp.Notification("can't set new deeplink")
		return
	}

	go func() {
		for url := range deeplink.Handler() {
			if _, err := robloxapp.Launch(url); err != nil {
				macosapp.Notification("can't launch roblox: " + err.Error())
			}
		}
	}()

	macosapp.Run(func() {
		urlhandler.Set(urlhandler.ROBLOX_BUNDLE_IDENTIFIER, "roblox-player")
	})

	// cc := make(chan io.Closer)
	// func(cc chan<- io.Closer) {
	// 	go func() {
	// 		for url := range deeplink.Handler() {
	// 			copy, err := robloxapp.NewCopy()
	// 			if err != nil {
	// 				macosapp.Notification("can't make new app copy: " + err.Error())
	// 				continue
	// 			}

	// 			go func() { // not the best solution, but it's alright =). cost of 1 goroutine is ~8kb ram.
	// 				cc <- copy
	// 			}()

	// 			syncbreaker.Break()
	// 			err = exec.Command("open", "-a", copy.Path(), url).Run()
	// 			if err != nil {
	// 				macosapp.Notification("can't open roblox app copy: " + err.Error())
	// 				continue
	// 			}

	// 			err = infoplist.SetMultipleInstancesProhibition(copy.Path()+"/Contents/Info.plist", false)
	// 			if err != nil {
	// 				macosapp.Notification("can't set multi instance prohibition: " + err.Error())
	// 				continue
	// 			}
	// 			syncbreaker.Break()
	// 		}
	// 	}()
	// }(cc)

	// macosapp.Run(func() { // handle app termination here
	// 	urlhandler.Set(urlhandler.ROBLOX_BUNDLE_IDENTIFIER, "roblox-player")

	// 	close(cc)
	// 	for c := range cc {
	// 		c.Close()
	// 	}
	// })
}
