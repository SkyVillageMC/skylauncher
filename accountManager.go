package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	client = http.Client{}
)

func Login(username, password string) error {
	payload := map[string]interface{}{
		"username":    username,
		"clientToken": "a921da84-13ad-4678-9b81-59e660116fcf",
		"requestUser": true,
	}

	body, _ := json.Marshal(payload)

	resp, err := client.Post(AuthorizeUrl, "application/json", strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(data))
	}

	CurrentSettings.Username = username
	CurrentSettings.IsLoggedIn = true
	return nil
}

func Logout() {
	CurrentSettings.Username = ""
	CurrentSettings.IsLoggedIn = false
}

func IsLoggedIn() bool {
	return GetSettings().IsLoggedIn
}

func GetUsername() string {
	return CurrentSettings.Username
}
