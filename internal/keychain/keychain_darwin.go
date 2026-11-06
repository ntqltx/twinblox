package keychain

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Foundation -framework Security
#include <stdlib.h>
#include "keychain.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/ntqltx/twinblox/internal/appdir"
)

var ErrNotFound = errors.New("keychain item not found")

// set function saves secret under account, replacing whatever was saved there before
func Set(account string, secret []byte) error {
	data := C.CBytes(secret) // copy it into memory so the objc can read it
	defer C.free(data)

	return call("write", account, func(service, account *C.char) C.OSStatus {
		return C.kcSet(service, account, data, C.int(len(secret)))
	})
}

// get function returns secret that was saved under given account
func Get(account string) ([]byte, error) {
	var data unsafe.Pointer
	var length C.int

	err := call("read", account, func(service, account *C.char) C.OSStatus {
		return C.kcGet(service, account, &data, &length)
	})

	if err != nil {
		return nil, err
	}
	defer C.free(data)
	return C.GoBytes(data, length), nil
}

// delete function basically just deletes secret that was saved under account
func Delete(account string) error {
	err := call("delete", account, func(service, account *C.char) C.OSStatus {
		return C.kcDelete(service, account)
	})
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	return err
}

// call turns the app's bundle id and the account name into C strings
func call(op, account string, fn func(service, account *C.char) C.OSStatus) error {
	s := C.CString(appdir.BUNDLE_ID)
	defer C.free(unsafe.Pointer(s))

	a := C.CString(account)
	defer C.free(unsafe.Pointer(a))

	switch status := fn(s, a); status {
	case C.errSecSuccess:
		return nil
	// checks keychain status code into an error
	case C.errSecItemNotFound:
		return ErrNotFound
	default:
		return fmt.Errorf("keychain %s failed: OSStatus %d", op, int(status))
	}
}
