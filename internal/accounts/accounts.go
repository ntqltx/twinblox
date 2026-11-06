package accounts

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/ntqltx/twinblox/internal/appdir"
	"github.com/ntqltx/twinblox/internal/keychain"
)

// account struct that can get callers
// without any private information about an account
type Account struct {
	Name     string `json:"name"`
	UserId   int64  `json:"userId"`
	Username string `json:"username"`
}

// returns list with every saved account
func List() ([]Account, error) {
	var list []Account
	err := appdir.Locked(func(dir string) error {
		var err error
		list, err = load(dir)
		return err
	})

	if list == nil {
		list = []Account{}
	}
	return list, err
}

func Save(acc Account) {
	// todo
}

// gets an account cookie SAFELY !!! using macos keychain
func Cookie(acc Account) ([]byte, error) {
	return keychain.Get(keychainKey(acc.UserId))
}

// returns saved account from the list under given name
func Find(name string) (Account, error) {
	list, err := List()
	if err != nil {
		return Account{}, err
	}

	for _, a := range list {
		if a.Name == name {
			return a, nil
		}
	}
	return Account{}, fmt.Errorf(
		"%w: %q",
		errors.New("no account with that name"),
		name,
	)
}

func load(dir string) ([]Account, error) {
	var list []Account
	err := appdir.ReadJSON(path(dir), &list)
	return list, err
}

func path(dir string) string {
	return filepath.Join(dir, "accounts.json")
}

func keychainKey(userId int64) string {
	return "roblosecurity:" + strconv.FormatInt(userId, 10)
}
