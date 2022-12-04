package main

import (
	"fmt"

	"go-chat/client/logger"
	"go-chat/client/model"
	"go-chat/client/process"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

//Login UI content
func LoginUI(Index chan int, userbox, groupbox, p2pbox *fyne.Container) (content *fyne.Container) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("input name")
	nameEntry.OnChanged = func(content string) {
		fmt.Println("name:", nameEntry.Text, "entered")
	}
	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("input password")

	nameBox := container.NewGridWithColumns(2, widget.NewLabel("Name"), nameEntry)
	passwordBox := container.NewGridWithColumns(2, widget.NewLabel("Password"), passEntry)

	statusLabel := widget.NewLabel("")

	SignUpBtn := widget.NewButton("SignUp", func() {
		Index <- 2
	})
	loginBtn := widget.NewButton("Login", func() {
		fmt.Println("name:", nameEntry.Text, "password:", passEntry.Text, "login")
		//cmd client return some status
		//use a func to change UI according to status
		up := process.UserProcess{}
		err := up.APPLogin(nameEntry.Text, passEntry.Text, userbox, groupbox, p2pbox)
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
func SignUpUI(Index chan int, userbox, groupbox, p2pbox *fyne.Container) (content *fyne.Container) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("input name")
	nameEntry.OnChanged = func(content string) {
		fmt.Println("name:", nameEntry.Text, "entered")
	}
	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("input password")
	passConfirmEntry := widget.NewPasswordEntry()
	passConfirmEntry.SetPlaceHolder("confirm password")

	nameBox := container.NewGridWithColumns(2, widget.NewLabel("Name"), nameEntry)
	passwordBox := container.NewGridWithColumns(2, widget.NewLabel("Password"), passEntry)
	passConfirmBox := container.NewGridWithColumns(2, widget.NewLabel("Password"), passConfirmEntry)

	statusLabel := widget.NewLabel("")

	SignUpBtn := widget.NewButton("Sign Up", func() {
		fmt.Println("name:", nameEntry.Text, "password:", passEntry.Text, "password confirm:", passConfirmEntry.Text, "SignUp")
		up := process.UserProcess{}
		err := up.APPRegister(nameEntry.Text, passEntry.Text, passConfirmEntry.Text, userbox, groupbox, p2pbox)
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

func MainUI(Index chan int) (content *container.AppTabs, userbox, groupbox, p2pbox *fyne.Container) {
	onlineUserBox := container.NewVBox()
	groupChatBox := container.NewVBox()
	P2PChatBox := container.NewVBox()

	userbox, groupbox, p2pbox = onlineUserBox, groupChatBox, P2PChatBox

	onlineUserScroll := container.NewVScroll(onlineUserBox)
	showOnlineUserBtn := widget.NewButton("Show online", func() {
		messageProcess := process.MessageProcess{}
		messageProcess.APPGetOnlineUerList(onlineUserBox, groupChatBox, P2PChatBox)
	})
	showOnlineUser := container.NewGridWithRows(2, onlineUserScroll, showOnlineUserBtn)

	groupChatScroll := container.NewVScroll(groupChatBox)
	groupChatEntry := widget.NewMultiLineEntry()
	groupChatEntry.SetPlaceHolder("Please enter something...")
	groupSendBtn := widget.NewButton("Send", func() {
		currentUser := model.CurrentUser
		messageProcess := process.MessageProcess{}
		text := widget.NewLabel(currentUser.UserName + " say to everyone: " + groupChatEntry.Text)
		groupChatBox.Add(text)
		err := messageProcess.SendGroupMessageToServer(0, currentUser.UserName, groupChatEntry.Text)
		if err != nil {
			logger.Error("Some error occurred when send data to server: %v\n", err)
		} else {
			logger.Success("Send group message succeed!\n\n")
		}
	})

	GroupChat := container.NewGridWithRows(3, groupChatScroll, groupChatEntry, groupSendBtn)

	P2PChatScroll := container.NewVScroll(P2PChatBox)
	P2PRecverEntry := widget.NewEntry()
	P2PRecverEntry.SetPlaceHolder("Please enter the receiver's name...")
	P2PChatEntry := widget.NewMultiLineEntry()
	P2PChatEntry.SetPlaceHolder("Please enter something you want to say...")
	P2PChatArea := container.NewHSplit(P2PRecverEntry, P2PChatEntry)
	P2PSendBtn := widget.NewButton("Send", func() {
		currentUser := model.CurrentUser
		text := widget.NewLabel(currentUser.UserName + " say to " + P2PRecverEntry.Text + " : " + P2PChatEntry.Text)
		P2PChatBox.Add(text)
		messageProcess := process.MessageProcess{}
		conn, err := messageProcess.PointToPointCommunication(P2PRecverEntry.Text, currentUser.UserName, P2PChatEntry.Text)
		if err != nil {
			logger.Error("Some error occurred when point to point comunication: %v\n", err)
			return
		}
		errMsg := make(chan error)
		go process.Response(conn, errMsg)
		err = <-errMsg

		if err.Error() != "<nil>" {
			logger.Error("Send message error: %v\n", err)
		}
	})
	P2PChat := container.NewGridWithRows(3, P2PChatScroll, P2PChatArea, P2PSendBtn)

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

	content, userbox, groupbox, p2pbox := MainUI(Index)
	w3.SetContent(content)
	w3.Resize(fyne.Size{Width: 800, Height: 600})
	w3.CenterOnScreen()

	w1.SetContent(LoginUI(Index, userbox, groupbox, p2pbox))
	w1.Resize(fyne.Size{Width: 500, Height: 300})
	w1.CenterOnScreen()

	w2.SetContent(SignUpUI(Index, userbox, groupbox, p2pbox))
	w2.Resize(fyne.Size{Width: 500, Height: 300})
	w2.CenterOnScreen()

	go changeWindow(&w1, &w2, &w3, Index)

	w1.Show()
	a.Run()
}
