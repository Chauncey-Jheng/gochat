package process

import (
	"encoding/json"
	"errors"
	"go-chat/client/logger"
	"go-chat/client/utils"
	common "go-chat/common/message"
	"go-chat/config"
	"net"

	"fyne.io/fyne/v2"
)

type UserProcess struct{}

// 用户登陆
func (up UserProcess) Login(userName, password string) (err error) {
	// connect server
	serverInfo := config.Configuration.ServerInfo
	conn, err := net.Dial("tcp", serverInfo.Host)

	if err != nil {
		logger.Error("Connect server error: %v", err)
		return
	}

	var message common.Message
	message.Type = common.LoginMessageType
	// 生成 loginMessage
	var loginMessage common.LoginMessage
	loginMessage.UserName = userName
	loginMessage.Password = password

	// 先序列话需要传到服务器的数据
	data, err := json.Marshal(loginMessage)
	if err != nil {
		logger.Error("Some error occurred when parse you data, error: %v\n", err)
		return
	}

	// 首先发送数据 data 的长度到服务器端
	// 将一个字符串的长度转为一个表示长度的切片
	message.Data = string(data)
	message.Type = common.LoginMessageType
	data, _ = json.Marshal(message)

	dispatcher := utils.Dispatcher{Conn: conn}
	err = dispatcher.SendData(data)
	if err != nil {
		logger.Error("Some error occurred when dispatch your data, error: %v\n", err)
		return
	}

	errMsg := make(chan error)
	go Response(conn, errMsg)
	err = <-errMsg

	if err != nil {
		logger.Error("Some error occurred with reponse, error: %v\n", err)
		return
	}

	return
}

// 用户登陆
func (up UserProcess) APPLogin(userName, password string, userbox, groupbox, p2pbox *fyne.Container) (err error) {
	// connect server
	serverInfo := config.Configuration.ServerInfo
	conn, err := net.Dial("tcp", serverInfo.Host)

	if err != nil {
		logger.Error("Connect server error: %v", err)
		return
	}

	var message common.Message
	message.Type = common.LoginMessageType
	// 生成 loginMessage
	var loginMessage common.LoginMessage
	loginMessage.UserName = userName
	loginMessage.Password = password

	// 先序列话需要传到服务器的数据
	data, err := json.Marshal(loginMessage)
	if err != nil {
		logger.Error("Some error occurred when parse you data, error: %v\n", err)
		return
	}

	// 首先发送数据 data 的长度到服务器端
	// 将一个字符串的长度转为一个表示长度的切片
	message.Data = string(data)
	message.Type = common.LoginMessageType
	data, _ = json.Marshal(message)

	dispatcher := utils.Dispatcher{Conn: conn}
	err = dispatcher.SendData(data)
	if err != nil {
		logger.Error("Some error occurred when dispatch your data, error: %v\n", err)
		return
	}

	errMsg := make(chan error)
	go APPResponse(conn, errMsg, userbox, groupbox, p2pbox)
	err = <-errMsg

	if err != nil {
		logger.Error("Some error occurred with reponse, error: %v\n", err)
		return
	}
	return
}

// 处理用户注册
func (up UserProcess) Register(userName, password, passwordConfirm string) (err error) {
	if password != passwordConfirm {
		err = errors.New("confirm password not match")
		return
	}
	serverInfo := config.Configuration.ServerInfo
	conn, err := net.Dial("tcp", serverInfo.Host)

	if err != nil {
		logger.Error("Connect server error: %v", err)
		return
	}

	// 定义消息
	var message common.Message

	// 生成 registerMessage
	var registerMessage common.RegisterMessage
	registerMessage.UserName = userName
	registerMessage.Password = password
	registerMessage.PasswordConfirm = passwordConfirm

	data, err := json.Marshal(registerMessage)
	if err != nil {
		logger.Error("Client occurred some error: %v\n", err)
	}

	// 构造需要传递给服务器的数据
	message.Data = string(data)
	message.Type = common.RegisterMessageType

	data, err = json.Marshal(message)
	if err != nil {
		logger.Error("RegisterMessage json Marshal error: %v\n", err)
		return
	}

	dispatcher := utils.Dispatcher{Conn: conn}
	err = dispatcher.SendData(data)
	if err != nil {
		logger.Error("Send data error!\n")
		return
	}

	errMsg := make(chan error)
	go Response(conn, errMsg)
	err = <-errMsg

	return
}

// 处理用户注册
func (up UserProcess) APPRegister(userName, password, passwordConfirm string, userbox, groupbox, p2pbox *fyne.Container) (err error) {
	if password != passwordConfirm {
		err = errors.New("confirm password not match")
		return
	}
	serverInfo := config.Configuration.ServerInfo
	conn, err := net.Dial("tcp", serverInfo.Host)

	if err != nil {
		logger.Error("Connect server error: %v", err)
		return
	}

	// 定义消息
	var message common.Message

	// 生成 registerMessage
	var registerMessage common.RegisterMessage
	registerMessage.UserName = userName
	registerMessage.Password = password
	registerMessage.PasswordConfirm = passwordConfirm

	data, err := json.Marshal(registerMessage)
	if err != nil {
		logger.Error("Client occurred some error: %v\n", err)
	}

	// 构造需要传递给服务器的数据
	message.Data = string(data)
	message.Type = common.RegisterMessageType

	data, err = json.Marshal(message)
	if err != nil {
		logger.Error("RegisterMessage json Marshal error: %v\n", err)
		return
	}

	dispatcher := utils.Dispatcher{Conn: conn}
	err = dispatcher.SendData(data)
	if err != nil {
		logger.Error("Send data error!\n")
		return
	}

	errMsg := make(chan error)
	go APPResponse(conn, errMsg, userbox, groupbox, p2pbox)
	err = <-errMsg

	return
}
