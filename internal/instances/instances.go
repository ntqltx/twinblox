package instances

import (
	"path/filepath"
	"time"

	"github.com/ntqltx/twinblox/internal/appdir"
	"github.com/ntqltx/twinblox/pkg/ps"
)

type Instance struct {
	Name      string    `json:"name"`
	UserId    int64     `json:"userId"`
	PID       int       `json:"pid"`
	PlaceId   int64     `json:"placeId"`
	StartedAt time.Time `json:"startedAt"`
}

func path(dir string) string {
	return filepath.Join(dir, "instances.json")
}

func alive(pid int) bool {
	p, err := ps.FindProcess(pid)
	return err == nil && p != nil && p.Executable() == "RobloxPlayer"
}

func LoadLocked(dir string) ([]Instance, error) {
	var list []Instance
	if err := appdir.ReadJSON(path(dir), &list); err != nil {
		return nil, err
	}

	live := []Instance{}
	for _, in := range list {
		if alive(in.PID) {
			live = append(live, in)
		}
	}

	if len(live) != len(list) {
		if err := appdir.WriteJSON(path(dir), live); err != nil {
			return nil, err
		}
	}

	return live, nil
}

func AddLocked(dir string, in Instance) error {
	list, err := LoadLocked(dir)
	if err != nil {
		return err
	}
	return appdir.WriteJSON(path(dir), append(list, in))
}

func List() ([]Instance, error) {
	var list []Instance

	err := appdir.Locked(func(dir string) error {
		var err error
		list, err = LoadLocked(dir)
		return err
	})

	return list, err
}
