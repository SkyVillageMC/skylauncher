package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"log"
	"skyvillage-launcher-rewrite/settings"
	"skyvillage-launcher-rewrite/updater"
	"skyvillage-launcher-rewrite/utils"
	"skyvillage-launcher-rewrite/views"
)

func main() {
	log.SetPrefix("SKYVILLAGE|> ")
	log.Printf("Starting launcher version %s by %s\n", utils.LauncherVersion, utils.Author)
	a := app.NewWithID("hu.skyvillage.launcher")

	utils.GameDir = utils.GetGameDir()

	settings.LoadSettings()

	views.SetOnLogin(func() {
		err := updater.CheckForUpdates()
		if err != nil {
			dialog.ShowError(err, views.MainWindow)
		}
	})

	views.InitMainWindow(a, func() {
		log.Println("Closing...")
		settings.Save()
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
	defer views.CloseMainWindow()
}
