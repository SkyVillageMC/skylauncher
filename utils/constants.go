package utils

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"fyne.io/fyne/v2"
)

var (
	LauncherVersion = "v0.0.1-pre"
	Author          = "Bendi"
	BaseUrl         = "https://bendi.tk/assets"
	ManifestUrl     = fmt.Sprintf("%s/version.json", BaseUrl)
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
		return "windows"

	case "darwin":
		return "osx"

	case "linux":
		return "linux"
	}
	return "not supported"
}

func IsLibraryCompatible(onlyIn string) bool {
	if onlyIn == "all" || onlyIn == "" {
		return true
	}
	return internalUpperOs() == onlyIn
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
