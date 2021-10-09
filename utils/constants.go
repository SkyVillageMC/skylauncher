package utils

import (
	"fmt"
	"fyne.io/fyne/v2"
	"log"
	"os"
	"path"
	"runtime"
)

var (
	LauncherVersion = "v0.0.1-pre"
	Author          = "Bendi"
	BaseUrl         = "https://bendimester23.tk/assets"
	ManifestUrl     = fmt.Sprintf("%s/1.14.4.json", BaseUrl)
	GameDir         = GetGameDir()
)

func GetGameDir() string {
	switch runtime.GOOS {
	case "linux":
		return path.Join(os.Getenv("HOME"), ".skyvillage")
	//case "ios":
	//case "darwin":
	//	return ""
	case "windows":
		return path.Join(os.Getenv("APPDATA"), ".skyvillage")
	default:
		fyne.NewNotification("Hiba!", "Nem támogatott operációs rendszer")
		log.Fatalf("Unsupported system!\n")
	}
	return ""
}

func internalUpperOs() string {
	switch runtime.GOOS {
	case "windows":
		return "WIN"

	case "darwin":
		return "OSX"

	case "linux":
		return "LINUX"
	}
	return "not supported"
}

func IsLibraryCompatible(onlyIn []string) bool {
	if len(onlyIn) == 0 {
		return true
	}
	for i := range onlyIn {
		if internalUpperOs() == onlyIn[i] {
			return true
		}
	}

	return false
}

func GetNativesName() string {
	switch runtime.GOOS {
	case "windows":
		return "natives-windows"

	case "darwin":
		return "natives-macos"

	case "linux":
		return "natives-linux"
	}
	return ""
}

func NeedMod(arr []Mod, elem string) bool {
	for _, i := range arr {
		if i.FileName == elem {
			return true
		}
	}
	return false
}
