package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"skyvillage-launcher-rewrite/auth"
)

var (
	mainWindow fyne.Window
)

func InitMainWindow(app fyne.App, close func()) {
	mainWindow = app.Driver().CreateWindow("SkyVillage Launcher")
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
	mainWindow.SetOnClosed(close)
	initContent(&mainWindow)
	mainWindow.ShowAndRun()
}

func CloseMainWindow() {
	mainWindow.Close()
}

func initContent(ww *fyne.Window) {
	w := *ww
	if auth.IsLoggedIn() {
		w.SetContent(getLoggedInContent())
	} else {
		w.SetContent(getLoginFormContent())
	}
}

func getLoggedInContent() fyne.CanvasObject {
	return container.NewVBox(
		getLoggedInNewsContent(),
		getLoggedInMenuBarContent(),
	)
}

func getLoggedInMenuBarContent() fyne.CanvasObject {
	return container.NewHBox(
		widget.NewLabel("Menu"),
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
	passwordInput.SetPlaceHolder("Nem kell megadni")
	passwordInput.Resize(fyne.NewSize(150, 30))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Felhasználónév", Widget: usernameInput},
			{Text: "Jelszó", Widget: passwordInput},
		},
		OnSubmit: func() {
			processLogin(usernameInput.Text, passwordInput.Text)
		},
		SubmitText: "Bejelentkezés",
	}

	return container.New(layout.NewCenterLayout(), widget.NewCard("               Bejelentkezés               ", " ", form))
}

func processLogin(username, password string) {
	auth.Login(username, password)
	initContent(&mainWindow)
}
