package main

import (
	"log"
	"skyvillage-launcher-rewrite/utils"

	"fyne.io/fyne/v2/app"
)

func main() {
	log.Printf("Starting launcher version %s.\n", utils.LauncherVersion)
	svApp = app.NewWithID("hu.skyvillage.launcher")

	utils.GameDir = utils.GetGameDir()

	LoadSettings()

	InitView()

	/* views.SetOnLogin(func() {
		err := updater.CheckForUpdates()
		if err != nil {
			dialog.ShowError(err, views.MainWindow)
		}
	}) */

	/* views.InitMainWindow(a, func() {
		log.Println("Closing...")
		SaveSettings()
	}, func() {
		if updater.CheckForLauncherUpdates() {
			err := updater.UpdateLauncher()
			if err != nil {
				dialog.ShowError(err, views.MainWindow)
				a.SendNotification(fyne.NewNotification("Hiba!", "Hiba történt a kliens frissítése közben!"))
				log.Fatalf("Error updating launcher!\n%s", err.Error())
			}
		}

	})
	defer views.CloseMainWindow() */
}
