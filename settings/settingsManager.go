package settings

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

var (
	CurrentSettings Settings
	defaultConfig   = Settings{
		Version:      1,
		OpenedBefore: false,
		Username:     "",
		MaxMemory:    1024,
		IsLoggedIn: false,
		DeveloperOptions: DeveloperOptions{
			IsDev:   false,
			DevCode: "NOPE",
		},
	}
)

type Settings struct {
	Version byte `json:"version"`
	OpenedBefore bool `json:"opened_before"`
	Username string `json:"username"`
	MaxMemory int `json:"max_memory"`
	IsLoggedIn bool `json:"is_logged_in"`
	DeveloperOptions DeveloperOptions `json:"developer_options"`
}

type DeveloperOptions struct {
	IsDev bool `json:"is_dev"`
	DevCode string `json:"dev_code"`
}

func LoadSettings() {
	log.Println("Loading CurrentSettings...")
	_, err := os.Stat(GetConfigPath())
	if err != nil {
		initDefaultConfig()
	}
	d, _ := ioutil.ReadFile(GetConfigPath())
	err = json.Unmarshal(d, &CurrentSettings)
	if err != nil {
		initDefaultConfig()
	}
}

func initDefaultConfig() {
	f, err := os.Create(GetConfigPath())
	if err != nil {
		log.Fatalf("Error creating config\n%s\n", err.Error())
	}
	err = os.MkdirAll(GetConfigPath(), os.ModeDir)
	if err != nil {
		log.Fatalf("Error creating config\n%s\n", err.Error())
	}
	d, _ := json.Marshal(defaultConfig)
	_, err = f.Write(d)
	if err != nil {
		log.Fatalf("Error creating config\n%s\n", err.Error())
	}
}

func GetConfigPath() string {
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("APPDATA") + "\\.skyvillage\\launcher.json"
	default:
		return os.Getenv("HOME") + "/.skyvillage/launcher.json"
	}
}

func Save() {
	log.Println("Saving CurrentSettings...")
	f, err := os.Create(GetConfigPath())
	if err != nil {
		log.Fatalf("Error saving config\n%s\n", err.Error())
	}
	d, _ := json.Marshal(CurrentSettings)
	_, err = f.Write(d)
	if err != nil {
		log.Fatalf("Error saving config\n%s\n", err.Error())
	}
}

func GetSettings() *Settings {
	return &CurrentSettings
}
