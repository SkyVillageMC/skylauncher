package updater

import (
	"errors"
	"log"
)

func CheckForLauncherUpdates() bool {
	log.Println("Checking for updates...")

	//latestVersion, err := utils.GetStringFromWeb("https://bendimester23.tk/assets/launcher_version")
	//if err != nil {
	//	log.Printf("Failed to check for launcher updates!\n%s", err.Error())
	//}
	//if semver.Compare(latestVersion, utils.LauncherVersion) != 0 {
	//	log.Printf("Update found!\nCurrent version: %s\nNew version:%s", utils.LauncherVersion, latestVersion)
	//	return true
	//}
	//log.Println("No update found!")
	return false
}

func UpdateLauncher() error {
	return errors.New("izé")
}
