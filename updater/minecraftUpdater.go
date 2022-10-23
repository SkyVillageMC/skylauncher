package updater

import (
	"fmt"
	"log"
	"os"
	"path"
	"skyvillage-launcher-rewrite/launcher"
	"skyvillage-launcher-rewrite/utils"
)

var (
	version    utils.VersionManifest
	mods       utils.ContentIndex
	assetIndex utils.Assets

	TaskChangeCallback func(string)
	TaskDoneCallback   func()
)

func CheckForUpdates() error {
	err := utils.GetStringAsJson(utils.ManifestUrl, &version)
	if err != nil {
		return err
	}
	log.Printf("Found version %s-%s.", version.Id, version.ReleaseType)
	launcher.SetVersion(version)
	if !utils.ExistsAndValid(path.Join(utils.GameDir, "client.jar"), version.Client.Hash) {
		downloadClient()
	}
	checkLibraries()
	checkNatives()
	checkAssets()
	checkMods()
	TaskChangeCallback("Kész.")
	utils.ProgressChangeCallback(0)
	TaskDoneCallback()
	return nil
}

func checkMods() error {
	err := utils.GetStringAsJson("https://bendi.tk/assets/content.json", &mods)
	if err != nil {
		return err
	}
	if !utils.FileExists(path.Join(utils.GameDir, "mods")) {
		os.Mkdir(path.Join(utils.GameDir, "mods"), os.ModeDir|os.ModePerm)
	}
	utils.ForEachFile(path.Join(utils.GameDir, "mods"), func(name string) {
		if !utils.NeedMod(mods.Mods, name) {
			log.Printf("Found suspicous mod: %s, removing it.\n", name)
			err := os.Remove(path.Join(utils.GameDir, "mods", name))
			if err != nil {
				log.Println(err.Error())
			}
		}
	})
	for i, m := range mods.Mods {
		if TaskChangeCallback != nil {
			TaskChangeCallback(fmt.Sprintf("(%d/%d) Tartalom letöltése: %s", i, len(mods.Mods), m.Name))
		}
		mP := path.Join(utils.GameDir, "mods", m.FileName)
		if !utils.ExistsAndValid(mP, m.Hash) {
			log.Printf("Downloading %s to %s.\nUrl: %s.\n", m.Name, mP, m.Url)
			err := utils.StartDownloadJob(m.Url, mP, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type LocalAssets []string

func (l LocalAssets) Has(hash string) bool {
	for _, v := range l {
		if hash == v {
			return true
		}
	}

	return false
}

func CollectLocalAssets() LocalAssets {
	var res LocalAssets = LocalAssets{}

	rootPath := path.Join(utils.GameDir, "assets", "objects")

	files, err := os.ReadDir(rootPath)
	if err != nil {
		log.Println(err.Error())
		return res
	}

	for _, v := range files {
		if v.IsDir() {
			files, err = os.ReadDir(path.Join(rootPath, v.Name()))
			if err != nil {
				log.Println(err.Error())
				return res
			}
			for _, v2 := range files {
				if !v2.IsDir() {
					res = append(res, v2.Name())
				}
			}
		} else {
			res = append(res, v.Name())
		}
	}

	return res
}

func checkAssets() {
	os.MkdirAll(path.Join(utils.GameDir, "assets", "indexes"), os.ModeDir|os.ModePerm)
	os.MkdirAll(path.Join(utils.GameDir, "assets", "objects"), os.ModeDir|os.ModePerm)
	indexUpdate := false
	assetIndexPath := path.Join(utils.GameDir, "assets", "indexes", fmt.Sprintf("%s.json", version.Assets.Id))
	if !utils.ExistsAndValid(assetIndexPath, version.Assets.Hash) {
		indexUpdate = true
		err := utils.StartDownloadJob(version.Assets.Url, assetIndexPath, true)
		if err != nil {
			log.Println(err.Error())
		}
	}
	err := utils.LoadJsonToPointer(assetIndexPath, &assetIndex)
	if err != nil {
		log.Println(err.Error())
	}
	c := 0
	max := len(assetIndex.Objects)

	localAssets := CollectLocalAssets()

	for f, o := range assetIndex.Objects {
		if TaskChangeCallback != nil {
			TaskChangeCallback(fmt.Sprintf("(%d/%d) Erőforrás letöltése: %s.", c, max, f))
			utils.ProgressChangeCallback((float64(c) / float64(max)) * 100)
		}

		aPath := path.Join(utils.GameDir, "assets", "objects", o.Hash[0:2], o.Hash)

		needUpdate := false
		if indexUpdate {
			needUpdate = !utils.ExistsAndValid(aPath, o.Hash)
		} else {
			needUpdate = !localAssets.Has(o.Hash)
		}

		if needUpdate {
			err := utils.StartDownloadJob(fmt.Sprintf("https://resources.download.minecraft.net/%s/%s", o.Hash[0:2], o.Hash), aPath, false)
			if err != nil {
				log.Println(err.Error())
			}
		}
		c++
	}
}

func checkNatives() {
	i := 0
	for _, lib := range version.Libraries {
		if lib.IsNative && utils.IsLibraryCompatible(lib.Os) {
			if checkNative(i, lib.Downloads["artifact"]) {
				if TaskChangeCallback != nil {
					TaskChangeCallback(fmt.Sprintf("(%d/15) Natív könyvtár frissítése: %s", i+1, lib.Name))
				}
				err := utils.StartDownloadJob(lib.Downloads["artifact"].Url, path.Join(utils.GameDir, "natives", fmt.Sprintf("%d.jar", i)), true)
				if err != nil {
					log.Println(err.Error())
				}
			}
			i++
		}
	}
}

func checkNative(i int, download utils.Download) bool {
	return !utils.ExistsAndValid(path.Join(utils.GameDir, "natives", fmt.Sprintf("%d.jar", i)), download.Hash)
}

func checkLibraries() {
	for i, lib := range version.Libraries {
		if checkLibrary(lib) {
			err := updateLibrary(lib, i)
			if err != nil {
				log.Printf("Error updating %s: %s\n", lib.Name, err.Error())
			}
		}
	}
}

func updateLibrary(lib utils.Library, index int) error {
	log.Printf("Updating library %s\n", lib.Name)
	if TaskChangeCallback != nil {
		TaskChangeCallback(fmt.Sprintf("(%d/%d) Könyvtár frissítése: %s", index, len(version.Libraries), lib.Name))
	}

	var err error = nil
	if lib.IsNative {
		err = utils.StartDownloadJob(lib.Downloads["artifact"].Url, path.Join(utils.GameDir, "natives", fmt.Sprintf("%d.jar", index)), true)
	} else {
		err = utils.StartDownloadJob(lib.Downloads["artifact"].Url, path.Join(utils.GameDir, "libraries", lib.Downloads["artifact"].Path), true)
	}
	return err
}

func checkLibrary(lib utils.Library) bool {
	if !utils.IsLibraryCompatible(lib.Os) {
		return false
	}
	return !utils.ExistsAndValid(path.Join(utils.GameDir, "libraries", lib.Downloads["artifact"].Path), lib.Downloads["artifact"].Hash)
}

func downloadClient() {
	if TaskChangeCallback != nil {
		TaskChangeCallback("Játék letöltése.")
	}
	err := utils.StartDownloadJob(version.Client.Url, path.Join(utils.GameDir, "client.jar"), true)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
}
