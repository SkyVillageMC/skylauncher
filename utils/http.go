package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var client = http.Client{
	Timeout: time.Second * 10,
}

var curentCounter *WriteCounter

var fn func(p float64)

type WriteCounter struct {
	Total float64
	Max   float64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += float64(n)
	fn(float64(wc.Total/wc.Max) * 100)
	return n, nil
}

func SetOnPrecChange(f func(p float64)) {
	fn = f
}

func GetStringAsJson(url string, target interface{}) error {
	raw, err := GetStringFromWeb(url)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(raw), target)
	if err != nil {
		return err
	}
	return nil
}

func GetStringFromWeb(url string) (string, error) {
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("SkyVillage Kliens %s", LauncherVersion))

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func StartDownloadJob(url, path string) error {
	if !FileExists(filepath.Dir(path)) {
		err := os.MkdirAll(filepath.Dir(path), 0777)
		if err != nil {
			return err
		}
	}
	out, err := os.Create(path + ".tmp")
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		err := out.Close()
		if err != nil {
			return err
		}
		return err
	}
	defer resp.Body.Close()

	counter := &WriteCounter{}
	counter.Max = float64(resp.ContentLength)

	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		err := out.Close()
		if err != nil {
			return err
		}
		return err
	}

	err = out.Close()
	if err != nil {
		return err
	}

	if err = os.Rename(path+".tmp", path); err != nil {
		return err
	}
	return nil
}
