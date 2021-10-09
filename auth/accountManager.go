package auth

import "skyvillage-launcher-rewrite/settings"

func Login(username, password string) {
	settings.CurrentSettings.Username = username
	settings.CurrentSettings.IsLoggedIn = true
}

func Logout() {
	settings.CurrentSettings.Username = ""
	settings.CurrentSettings.IsLoggedIn = false
}

func IsLoggedIn() bool {
	return settings.GetSettings().IsLoggedIn
}

func GetUsername() string {
	return settings.CurrentSettings.Username
}
