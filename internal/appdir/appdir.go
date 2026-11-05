package appdir

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

const NAME string = "Twinblox"
const BUNDLE_ID string = "com.ntqltx.twinblox"

func Dir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	d := filepath.Join(base, NAME)
	if err := os.MkdirAll(d, 0o700); err != nil {
		return "", err
	}

	return d, nil
}

func Locked(fn func(dir string) error) error {
	d, err := Dir()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join(d, ".lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return err
	}

	defer f.Close()
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}

	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return fn(d)
}

func ReadJSON(path string, v any) error {
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func WriteJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
