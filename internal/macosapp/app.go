package macosapp

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa
#include "app.h"
*/
import "C"

import (
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/ntqltx/twinblox/internal/instances"
)

var onOpenURL func(string)
var onQuit func()

func init() {
	runtime.LockOSThread()
}

func Run(openURL func(url string), quit func()) {
	onOpenURL = openURL
	onQuit = quit

	C.RunApp()
}

func isTrayProcess(pid int) bool {
	out, err := exec.Command("ps", "-o", "args=", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "-launchToTray")
}

//export goIsTrayProcess
func goIsTrayProcess(pid C.int) C.bool {
	return C.bool(isTrayProcess(int(pid)))
}

//export goInstanceLabel
func goInstanceLabel(pid C.int) *C.char {
	list, err := instances.List()
	if err != nil {
		return nil
	}

	for _, in := range list {
		if in.PID == int(pid) {
			return C.CString(in.Name)
		}
	}

	return nil
}

//export goOpenURL
func goOpenURL(url *C.char) {
	onOpenURL(C.GoString(url))
}

//export goQuit
func goQuit() {
	onQuit()
}
