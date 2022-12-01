package main

import (
	"bufio"
	"fmt"
	"go-chat/client/logger"
	"go-chat/client/model"
	"go-chat/client/process"
	"os"
)

func commandLine() {
	var (
		key             int
		loop            = true
		userName        string
		password        string
		passwordConfirm string
	)

	for loop {
		logger.Info("\n----------------Welcome to the chat room--------------\n")
		logger.Info("\t\tSelect the options:\n")
		logger.Info("\t\t\t 1、Sign in\n")
		logger.Info("\t\t\t 2、Sign up\n")
		logger.Info("\t\t\t 3、Exit the system\n")

		// get user input
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			logger.Info("sign In Please\r\n")
			logger.Notice("Username:\n")
			fmt.Scanf("%s\n", &userName)
			logger.Notice("Password:\n")
			fmt.Scanf("%s\n", &password)

			up := process.UserProcess{}
			err := up.Login(userName, password)

			if err != nil {
				logger.Error("Login failed: %v\r\n", err)
			} else {
				logger.Success("Login succeed!\r\n")
				for {
					showAfterLoginMenu()
				}
			}
		case 2:
			logger.Info("Create account\n")
			logger.Notice("user name:\n")
			fmt.Scanf("%s\n", &userName)
			logger.Notice("password:\n")
			fmt.Scanf("%s\n", &password)
			logger.Notice("password confirm:\n")
			fmt.Scanf("%s\n", &passwordConfirm)

			up := process.UserProcess{}
			err := up.Register(userName, password, passwordConfirm)
			if err != nil {
				logger.Error("Create account failed: %v\n", err)
			} else {
				logger.Success("SignUp succeed!\r\n")
			}
		case 3:
			logger.Warn("Exit...\n")
			loop = false // this is equal to 'os.Exit(0)'
		default:
			logger.Error("Select is invalid!\n")
		}
	}
}

// 登陆成功菜单显示：
func showAfterLoginMenu() {
	logger.Info("\n----------------login succeed!----------------\n")
	logger.Info("\t\tselect what you want to do\n")
	logger.Info("\t\t1. Show all online users\n")
	logger.Info("\t\t2. Send group message\n")
	logger.Info("\t\t3. Point-to-point communication\n")
	logger.Info("\t\t4. Exit\n")
	var key int
	var content string
	var inputReader *bufio.Reader
	var err error
	inputReader = bufio.NewReader(os.Stdin)

	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		messageProcess := process.MessageProcess{}
		err = messageProcess.GetOnlineUerList()
		if err != nil {
			logger.Error("Some error occurred when get online user list, error: %v\n", err)
		}
	case 2:
		logger.Notice("Say something:\n")
		content, err = inputReader.ReadString('\n')
		if err != nil {
			logger.Error("Some error occurred when you input, error: %v\n", err)
		}
		currentUser := model.CurrentUser
		messageProcess := process.MessageProcess{}
		err = messageProcess.SendGroupMessageToServer(0, currentUser.UserName, content)
		if err != nil {
			logger.Error("Some error occurred when send data to server: %v\n", err)
		} else {
			logger.Success("Send group message succeed!\n\n")
		}
	case 3:
		var targetUserName string

		logger.Notice("Select one friend by user name\n")
		fmt.Scanf("%s\n", &targetUserName)
		logger.Notice("Input message:\n")
		content, err = inputReader.ReadString('\n')
		if err != nil {
			logger.Error("Some error occurred when you input, error: %v\n", err)
		}
		messageProcess := process.MessageProcess{}
		conn, err := messageProcess.PointToPointCommunication(targetUserName, model.CurrentUser.UserName, content)
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
	case 4:
		logger.Warn("Exit...\n")
		os.Exit(0)
	default:
		logger.Info("Selected invalid!\n")
	}
}
