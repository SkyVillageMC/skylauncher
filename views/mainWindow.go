package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"skyvillage-launcher-rewrite/auth"
	"skyvillage-launcher-rewrite/launcher"
	"skyvillage-launcher-rewrite/utils"
)

var (
	MainWindow  fyne.Window
	progressBar *widget.ProgressBar
	taskLabel   *widget.Label
	onLogin     func()
	playBtn     *widget.Button
	a           *fyne.App
)

func InitMainWindow(app fyne.App, close func(), preShow func()) {
	a = &app
	MainWindow = app.Driver().CreateWindow("SkyVillage Launcher")
	MainWindow.Resize(fyne.Size{
		Height: 520,
		Width:  470,
	})

	icon, err := fyne.LoadResourceFromPath("./icon.png")
	if err == nil {
		MainWindow.SetIcon(icon)
	}

	MainWindow.CenterOnScreen()
	MainWindow.SetFixedSize(true)
	MainWindow.SetOnClosed(close)
	initContent(&MainWindow)
	go preShow()
	MainWindow.ShowAndRun()
}

func CloseMainWindow() {
	if MainWindow != nil {
		MainWindow.Close()
	}
}

func initContent(ww *fyne.Window) {
	progressBar = widget.NewProgressBar()
	progressBar.Min = 0
	progressBar.Max = 100
	utils.SetOnPrecChange(func(p float64) {
		progressBar.SetValue(p)
	})

	taskLabel = widget.NewLabel(" ")
	taskLabel.Wrapping = fyne.TextTruncate

	w := *ww
	if auth.IsLoggedIn() {
		w.SetContent(getLoggedInContent())
	} else {
		w.SetContent(getLoginFormContent())
	}
}

func SetProgressBar(val float64) {
	progressBar.SetValue(val)
	progressBar.Refresh()
}

func SetCurrentTask(task string) {
	taskLabel.SetText(task)
}
func SetOnLogin(handler func()) {
	onLogin = handler
}

func getLoggedInContent() fyne.CanvasObject {
	go onLogin()
	return container.NewVBox(
		getLoggedInNewsContent(),
		getLoggedInMenuBarContent(),
	)
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
		MainWindow.Close()
		MainWindow = nil
		a = nil
		err := launcher.LaunchGame(auth.GetUsername(), "1aab1940-9f93-4dc8-a234-b9961d915ce2", "asd")
		if err != nil {
			log.Println(err.Error())
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
					auth.Logout()
					initContent(&MainWindow)
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
		fyne.NewSize(460, 400),
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
	initContent(&MainWindow)
}
