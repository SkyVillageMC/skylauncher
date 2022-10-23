package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var client = http.Client{
	Timeout: time.Second * 5,
}

var (
	ProgressChangeCallback func(float64)
)

type WriteCounter struct {
	Total float64
	Max   float64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += float64(n)
	ProgressChangeCallback(float64(wc.Total/wc.Max) * 100)
	return n, nil
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
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func StartDownloadJob(url, path string, setProgress bool) error {
	if !FileExists(filepath.Dir(path)) {
		err := os.MkdirAll(filepath.Dir(path), 0777)
		if err != nil {
			return err
		}
	}
	if setProgress {
		ProgressChangeCallback(0)
	}
	out, err := os.Create(path + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		err := out.Close()
		if err != nil {
			return err
		}
		return err
	}
	defer resp.Body.Close()

	if setProgress {
		counter := &WriteCounter{}
		counter.Max = float64(resp.ContentLength)

		if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
			return err
		}
	} else {

		if _, err = io.Copy(out, resp.Body); err != nil {
			return err
		}
	}

	if err = os.Rename(path+".tmp", path); err != nil {
		return err
	}
	return nil
}
