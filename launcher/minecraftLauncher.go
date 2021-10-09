package launcher

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"skyvillage-launcher-rewrite/utils"
)

var (
	version utils.VersionManifest
)

func SetVersion(v utils.VersionManifest) {
	version = v
}

func LaunchGame(username, uuid, sessionId string) error {
	runtime.GC()
	cmd := exec.Command("java", GetArgs(username, uuid, sessionId)...)
	cmd.Stdout = os.Stdout
	cmd.Dir = utils.GameDir

	version = *new(utils.VersionManifest)
	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetArgs(username, uuid, sessionId string) []string {
	return []string{
		"-Xmx1G",
		"-cp",
		getCp(),
		version.MainClass,
		"--gameDir",
		utils.GameDir,
		"--assetsDir",
		path.Join(utils.GameDir, "assets"),
		"--assetIndex",
		version.AssetsVersion,
		"--uuid",
		uuid,
		"--accessToken",
		sessionId,
		"--version",
		version.Id,
		"--userType",
		"mojang",
		"--versionType",
		"release",
		"--username",
		username,
	}
}

func getCp() string {
	cp := ""
	os := runtime.GOOS
	if os != "linux" {
		cp = "\""
	}
	for _, lib := range version.Libraries {
		if utils.IsLibraryCompatible(lib.OnlyIn) {
			if os == "linux" {
				cp = fmt.Sprintf("%s%s:", cp, utils.GetLibraryPath(lib.Name))
			} else {
				cp = fmt.Sprintf("%s%s;", cp, utils.GetLibraryPath(lib.Name))
			}
		}
	}
	utils.ForEachFile(path.Join(utils.GameDir, "natives"), func(n string) {
		cp = fmt.Sprintf("%s%s:", cp, path.Join(utils.GameDir, "natives", n))
	})
	cp = fmt.Sprintf("%s%s", cp, path.Join(utils.GameDir, "client.jar"))
	if os != "linux" {
		cp = fmt.Sprintf("%s\"", cp)
	}
	return cp
}
