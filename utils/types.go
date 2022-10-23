package utils

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

type VersionManifest struct {
	Id          string     `json:"id,omitempty"`
	ReleaseType string     `json:"type,omitempty"`
	Libraries   []Library  `json:"libraries,omitempty"`
	MainClass   string     `json:"mainClass,omitempty"`
	Client      Download   `json:"client"`
	Assets      AssetIndex `json:"assetIndex"`
}

type Library struct {
	Name      string              `json:"name,omitempty"`
	IsNative  bool                `json:"isNative"`
	Os        string              `json:"os"`
	Downloads map[string]Download `json:"downloads,omitempty"`
}

type Download struct {
	Url  string `json:"url,omitempty"`
	Hash string `json:"sha1,omitempty"`
	Size int32  `json:"size,omitempty"`
	Path string `json:"path,omitempty"`
}

type AssetIndex struct {
	TotalSize int32  `json:"totalSize,omitempty"`
	Id        string `json:"id,omitempty"`
	Url       string `json:"url,omitempty"`
	Hash      string `json:"sha1,omitempty"`
	Size      int32  `json:"size,omitempty"`
}

type Assets struct {
	Objects map[string]Object `json:"objects"`
}

type Mod struct {
	Name     string `json:"name,omitempty"`
	FileName string `json:"file_name,omitempty"`
	Url      string `json:"url,omitempty"`
	Hash     string `json:"sha1,omitempty"`
}

type ContentIndex struct {
	Mods []Mod `json:"mods,omitempty"`
}

type Object struct {
	Hash string `json:"hash,omitempty"`
	Size int32  `json:"size,omitempty"`
}

var libDir = path.Join(GameDir, "libraries")

func GetLibraryPath(in string) string {
	i := strings.Split(in, ":")
	lPath := path.Join(strings.ReplaceAll(i[0], ".", string(filepath.Separator)), i[1], i[2])
	name := fmt.Sprintf("%s-%s.jar", i[1], i[2])
	return path.Join(libDir, lPath, name)
}
