package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func LoadJsonToPointer(path string, dest interface{}) error {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(raw, dest)
	return err
}

func ForEachFile(path string, f func(name string)) {
	dir, err := os.ReadDir(path)
	if err != nil {
		log.Println(err.Error())
	}
	for _, i := range dir {
		f(i.Name())
	}
}
