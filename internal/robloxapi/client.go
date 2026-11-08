package robloxapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

type User struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

// this function returns the user the session cookie belongs to
func Authenticated(cookie string) (User, error) {
	req, err := newRequest(http.MethodGet, "https://users.roblox.com/v1/users/authenticated", cookie)
	if err != nil {
		return User{}, nil
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return User{}, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return User{}, errors.New("roblox session expired, sign in to this account again")
	}

	if resp.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("users.roblox.com returned %s", resp.Status)
	}

	var u User
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return User{}, err
	}
	return u, nil
}

func newRequest(method, url, cookie string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, strings.NewReader("{}"))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", ".ROBLOSECURITY="+cookie)
	req.Header.Set("Referer", "https://www.roblox.com/")
	req.Header.Set("Origin", "https://www.roblox.com")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
