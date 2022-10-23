package main

import (
	"errors"
	"log"
	"skyvillage-launcher-rewrite/launcher"
	"skyvillage-launcher-rewrite/updater"
	"skyvillage-launcher-rewrite/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	mainWindow  fyne.Window
	svApp       fyne.App
	progressBar *widget.ProgressBar
	taskLabel   *widget.Label
	playBtn     *widget.Button
)

func getLoggedInContent() fyne.CanvasObject {
	progressBar = widget.NewProgressBar()
	progressBar.Min = 0
	progressBar.Max = 100

	taskLabel = widget.NewLabel(" ")
	taskLabel.Wrapping = fyne.TextTruncate

	utils.ProgressChangeCallback = SetProgressBar
	updater.TaskChangeCallback = SetCurrentTask
	updater.TaskDoneCallback = AllowPlay

	return container.NewVBox(
		getLoggedInNewsContent(),
		getLoggedInMenuBarContent(),
	)
}

func SetProgressBar(val float64) {
	progressBar.SetValue(val)
	progressBar.Refresh()
}

func SetCurrentTask(task string) {
	taskLabel.SetText(task)
	taskLabel.Refresh()
}

func AllowPlay() {
	playBtn.Enable()
	playBtn.Refresh()
}

func DisallowPlay() {
	playBtn.Disable()
	playBtn.Refresh()
}

func getLoggedInMenuBarContent() fyne.CanvasObject {
	playBtn = widget.NewButton("                                   Játék                                   ", func() {
		DisallowPlay()
		mainWindow.Hide()
		err := launcher.LaunchGame(GetUsername(), "1aab1940-9f93-4dc8-a234-b9961d915ce2", "asd")
		if err != nil {
			log.Println(err.Error())
			mainWindow.Show()
			dialog.ShowError(errors.New("A játék kicrashelt. :("), mainWindow)
			AllowPlay()
		} else {
			mainWindow.Close()
		}
	})
	playBtn.Importance = widget.HighImportance
	playBtn.DisableableWidget.Disable()
	return container.NewVBox(container.NewHBox(
		container.NewHBox(
			playBtn,
			container.NewVBox(
				widget.NewButton("Beállítások", func() {

				}),
				widget.NewButton("Kijelentkezés", func() {
					Logout()
					mainWindow.SetContent(getLoginFormContent())
				}),
			),
		),
	),
		progressBar,
		taskLabel,
	)
}

func getLoggedInNewsContent() fyne.CanvasObject {
	return container.NewGridWrap(
		fyne.NewSize(440, 400),
		container.NewCenter(
			widget.NewCard("Hírek", "",
				container.NewGridWrap(
					fyne.NewSize(400, 300),
					widget.NewLabel("Igen"),
				),
			),
		),
	)
}

func getLoginFormContent() fyne.CanvasObject {
	usernameInput := widget.NewEntry()
	usernameInput.Resize(fyne.NewSize(150, 30))
	usernameInput.Refresh()
	passwordInput := widget.NewPasswordEntry()
	passwordInput.Resize(fyne.NewSize(150, 30))

	var form *widget.Form
	loading := false
	form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Felhasználónév", Widget: usernameInput},
			{Text: "Jelszó", Widget: passwordInput},
		},
		OnSubmit: func() {
			if loading {
				return
			}
			loading = true
			form.SubmitText = "Bejeletkezés..."

			defer func() {
				loading = false
			}()

			if err := Login(usernameInput.Text, passwordInput.Text); err != nil {
				log.Println(err.Error())
				dialog.ShowError(errors.New("Hiba történt!"), mainWindow)
				return
			}
			mainWindow.SetContent(getLoggedInContent())
		},
		SubmitText: "Bejelentkezés",
	}

	return container.New(layout.NewCenterLayout(), widget.NewCard("               Bejelentkezés               ", " ", form))
}

func InitView() {
	mainWindow = svApp.NewWindow("SkyVillage Launcher")
	mainWindow.Resize(fyne.Size{
		Height: 500,
		Width:  450,
	})

	icon, err := fyne.LoadResourceFromPath("./icon.png")
	if err == nil {
		mainWindow.SetIcon(icon)
	}

	mainWindow.CenterOnScreen()
	mainWindow.SetFixedSize(true)

	if IsLoggedIn() {
		mainWindow.SetContent(getLoggedInContent())
	} else {
		mainWindow.SetContent(getLoginFormContent())
	}
	go updater.CheckForUpdates()
	defer SaveSettings()
	mainWindow.ShowAndRun()
}
