package updater

import (
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"log"
	"os"
	"path"
	"skyvillage-launcher-rewrite/launcher"
	"skyvillage-launcher-rewrite/utils"
	"skyvillage-launcher-rewrite/views"
)

var (
	version    utils.VersionManifest
	mods       utils.ContentIndex
	assetIndex utils.Assets
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
	views.SetProgressBar(0)
	views.SetCurrentTask("Indításra kész!")
	views.AllowPlay()
	return nil
}

func checkMods() {
	err := utils.GetStringAsJson("https://bendimester23.tk/assets/content.json", &mods)
	if err != nil {
		log.Println(err.Error())
		dialog.ShowError(err, views.MainWindow)
		return
	}
	utils.ForEachFile(path.Join(utils.GameDir, "mods"), func(name string) {
		if !utils.NeedMod(mods.Mods, name) {
			err := os.Remove(path.Join(utils.GameDir, "mods", name))
			if err != nil {
				log.Println(err.Error())
			}
		}
	})
	for i, m := range mods.Mods {
		views.SetCurrentTask(fmt.Sprintf("(%d/%d) Tartalom letöltése: %s", i, len(mods.Mods), m.Name))
		mP := path.Join(utils.GameDir, "mods", m.FileName)
		if !utils.ExistsAndValid(mP, m.Hash) {
			err := utils.StartDownloadJob(m.Url, mP)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

func checkAssets() {
	assetIndexPath := path.Join(utils.GameDir, "assets", "indexes", fmt.Sprintf("%s.json", version.AssetsVersion))
	if !utils.ExistsAndValid(assetIndexPath, version.Assets.Hash) {
		err := utils.StartDownloadJob(version.Assets.Url, assetIndexPath)
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
	for f, o := range assetIndex.Objects {
		views.SetCurrentTask(fmt.Sprintf("(%d/%d) Erőforrás letöltése: %s", c, max, f))
		aPath := path.Join(utils.GameDir, "assets", "objects", o.Hash[0:2], o.Hash)
		if !utils.ExistsAndValid(aPath, o.Hash) {
			err := utils.StartDownloadJob(fmt.Sprintf("https://resources.download.minecraft.net/%s/%s", o.Hash[0:2], o.Hash), aPath)
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
		if lib.Downloads[utils.GetNativesName()].Url != "" && utils.IsLibraryCompatible(lib.OnlyIn) {
			if checkNative(i, lib.Downloads[utils.GetNativesName()]) {
				views.SetCurrentTask(fmt.Sprintf("(%d/15) Natív könyvtár frissítése: %s", i+1, lib.Name))
				err := utils.StartDownloadJob(lib.Downloads[utils.GetNativesName()].Url, path.Join(utils.GameDir, "natives", fmt.Sprintf("%d.jar", i)))
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
			updateLibrary(lib, i)
		}
	}
}

func updateLibrary(lib utils.Library, index int) {
	log.Printf("Updating library %s\n", lib.Name)
	views.SetCurrentTask(fmt.Sprintf("(%d/%d) Könyvtár frissítése: %s", index, len(version.Libraries), lib.Name))
	err := utils.StartDownloadJob(lib.Downloads["artifact"].Url, utils.GetLibraryPath(lib.Name))
	if err != nil {
		log.Println(err.Error())
	}
	if lib.Downloads[utils.GetNativesName()].Url != "" {
		err := utils.StartDownloadJob(lib.Downloads[utils.GetNativesName()].Url, path.Join(utils.GameDir, "natives", fmt.Sprintf("%d.jar", index)))
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func checkLibrary(lib utils.Library) bool {
	if !utils.IsLibraryCompatible(lib.OnlyIn) {
		return false
	}
	return !utils.ExistsAndValid(utils.GetLibraryPath(lib.Name), lib.Downloads["artifact"].Hash)
}

func downloadClient() {
	views.SetCurrentTask("Játék letöltése.")
	err := utils.StartDownloadJob(version.Client.Url, path.Join(utils.GameDir, "client.jar"))
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
}
