package main

import (
	"fmt"

	"go-chat/client/logger"
	"go-chat/client/process"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

//Login UI content
func LoginUI(Index chan int) (content *fyne.Container) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("input name")
	nameEntry.OnChanged = func(content string) {
		fmt.Println("name:", nameEntry.Text, "entered")
	}
	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("input password")

	nameBox := container.NewHBox(widget.NewLabel("Name"), layout.NewSpacer(), nameEntry)
	passwordBox := container.NewHBox(widget.NewLabel("Password"), layout.NewSpacer(), passEntry)

	statusLabel := widget.NewLabel("")

	SignUpBtn := widget.NewButton("SignUp", func() {
		Index <- 2
	})
	loginBtn := widget.NewButton("Login", func() {
		fmt.Println("name:", nameEntry.Text, "password:", passEntry.Text, "login")
		//cmd client return some status
		//use a func to change UI according to status
		up := process.UserProcess{}
		err := up.Login(nameEntry.Text, passEntry.Text)
		if err != nil {
			s := fmt.Sprintf("Login failed: %v\r\n", err)
			logger.Error(s)
			statusLabel.SetText(s)
		} else {
			logger.Success("Login succeed.\r\n")
			Index <- 3
			fmt.Println("index <- 3 down")
		}
	})

	content = container.NewVBox(nameBox, passwordBox, loginBtn, SignUpBtn, statusLabel)
	return
}

//SignUp UI content
func SignUpUI(Index chan int) (content *fyne.Container) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("input name")
	nameEntry.OnChanged = func(content string) {
		fmt.Println("name:", nameEntry.Text, "entered")
	}
	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("input password")
	passConfirmEntry := widget.NewPasswordEntry()
	passConfirmEntry.SetPlaceHolder("confirm password")

	nameBox := container.NewHBox(widget.NewLabel("Name"), layout.NewSpacer(), nameEntry)
	passwordBox := container.NewHBox(widget.NewLabel("Password"), layout.NewSpacer(), passEntry)
	passConfirmBox := container.NewHBox(widget.NewLabel("Password"), layout.NewSpacer(), passConfirmEntry)

	statusLabel := widget.NewLabel("")

	SignUpBtn := widget.NewButton("Sign Up", func() {
		fmt.Println("name:", nameEntry.Text, "password:", passEntry.Text, "password confirm:", passConfirmEntry.Text, "SignUp")
		up := process.UserProcess{}
		err := up.Register(nameEntry.Text, passEntry.Text, passConfirmEntry.Text)
		if err != nil {
			s := fmt.Sprintf("Create account failed: %v\n", err)
			logger.Error(s)
			statusLabel.SetText(s)
		} else {
			logger.Success("SignUp succeed!\r\n")
			statusLabel.SetText("SignUp succeed!\r\n")
		}
	})

	loginBtn := widget.NewButton("Login", func() {
		Index <- 1
	})

	content = container.NewVBox(nameBox, passwordBox, passConfirmBox, SignUpBtn, loginBtn, statusLabel)
	return
}

func MainUI(Index chan int) (content *container.AppTabs) {
	showOnlineUser := container.NewVBox()
	GroupChat := container.NewVBox()
	P2PChat := container.NewVBox()

	content = container.NewAppTabs(
		container.NewTabItem("Show online user", showOnlineUser),
		container.NewTabItem("Group chat", GroupChat),
		container.NewTabItem("P2P chat", P2PChat),
	)

	content.SetTabLocation(container.TabLocationLeading)
	return
}

//Change UI
func changeWindow(w1 *fyne.Window, w2 *fyne.Window, w3 *fyne.Window, Index chan int) {
	for {
		w := <-Index
		switch w {
		case 1:
			fmt.Println("case 1")
			(*w2).Hide()
			(*w3).Hide()
			(*w1).Show()
		case 2:
			fmt.Println("case 2")
			(*w1).Hide()
			(*w3).Hide()
			(*w2).Show()
		case 3:
			fmt.Println("case 3")
			(*w1).Hide()
			(*w2).Hide()
			(*w3).Show()
		}
	}
}

func GUI() {
	a := app.New()
	w1 := a.NewWindow("login")
	w2 := a.NewWindow("SignUp")
	w3 := a.NewWindow("GOCHAT")
	Index := make(chan int)

	w1.SetContent(LoginUI(Index))
	w2.SetContent(SignUpUI(Index))
	w3.SetContent(MainUI(Index))

	go changeWindow(&w1, &w2, &w3, Index)

	w1.Show()
	a.Run()
}
