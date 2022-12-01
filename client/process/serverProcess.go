package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-chat/client/logger"
	"go-chat/client/model"
	"go-chat/client/utils"
	common "go-chat/common/message"
	"net"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func dealLoginResponse(responseMsg common.ResponseMessage) (err error) {
	switch responseMsg.Code {
	case 200:
		// 解析当前用户信息
		var userInfo common.UserInfo
		err = json.Unmarshal([]byte(responseMsg.Data), &userInfo)
		if err != nil {
			return
		}

		// 初始化 CurrentUser
		user := model.User{}
		err = user.InitCurrentUser(userInfo.ID, userInfo.UserName)
		logger.Success("Login succeed!\n")
		logger.Notice("Current user, id: %d, name: %v\n", model.CurrentUser.UserID, model.CurrentUser.UserName)
		if err != nil {
			return
		}
	case 500:
		err = errors.New("server error")
	case 404:
		err = errors.New("user does not exist")
	case 403:
		err = errors.New("password invalid")
	default:
		err = errors.New("some error")
	}
	return
}

func dealRegisterResponse(responseMsg common.ResponseMessage) (err error) {
	switch responseMsg.Code {
	case 200:
		logger.Success("Register succeed!\n")
	case 500:
		err = errors.New("server error")
	case 403:
		err = errors.New("user already exists")
	case 402:
		err = errors.New("password invalid")
	default:
		err = errors.New("some error")
	}
	return
}

func dealGroupMessage(responseMsg common.ResponseMessage) (err error) {
	var groupMessage common.SendGroupMessageToClient
	err = json.Unmarshal([]byte(responseMsg.Data), &groupMessage)
	if err != nil {
		return
	}
	logger.Info("%v send you:", groupMessage.UserName)
	logger.Notice("\t%v\n", groupMessage.Content)
	return
}

func APPdealGroupMessage(responseMsg common.ResponseMessage, box *fyne.Container) (err error) {
	var groupMessage common.SendGroupMessageToClient
	err = json.Unmarshal([]byte(responseMsg.Data), &groupMessage)
	if err != nil {
		return
	}
	logger.Info("%v send you (APP):", groupMessage.UserName)
	logger.Notice("\t%v\n", groupMessage.Content)

	s := fmt.Sprintf("%v send you: %v", groupMessage.UserName, groupMessage.Content)
	text := widget.NewLabel(s)
	text.Wrapping = fyne.TextWrapBreak
	box.Add(text)
	return
}

func showAllOnlineUsersList(responseMsg common.ResponseMessage) (err error) {
	if responseMsg.Code != 200 {
		err = errors.New("server Error")
		return
	}

	var userList []common.UserInfo
	err = json.Unmarshal([]byte(responseMsg.Data), &userList)
	if err != nil {
		return
	}

	logger.Success("Online user list(%v users)\n", len(userList))
	logger.Notice("\t\tID\t\tname\n")
	for _, info := range userList {
		logger.Success("\t\t%v\t\t%v\n", info.ID, info.UserName)
	}

	return
}

func APPshowAllOnlineUsersList(responseMsg common.ResponseMessage, box *fyne.Container) (err error) {
	if responseMsg.Code != 200 {
		err = errors.New("server Error")
		return
	}

	var userList []common.UserInfo
	err = json.Unmarshal([]byte(responseMsg.Data), &userList)
	if err != nil {
		return
	}

	logger.Success("Online user list(%v users)\n", len(userList))
	logger.Notice("\t\tID\t\tname\n")
	box.RemoveAll()
	for _, info := range userList {
		s := fmt.Sprintf("\t\t%v\t\t%v\n", info.ID, info.UserName)
		logger.Success(s)
		text := widget.NewLabel(s)
		box.Add(text)
	}
	return
}

func showPointToPointMessage(responseMsg common.ResponseMessage) (err error) {
	if responseMsg.Code != 200 {
		err = errors.New(responseMsg.Error)
		return
	}

	var pointToPointMessage common.PointToPointMessage
	err = json.Unmarshal([]byte(responseMsg.Data), &pointToPointMessage)
	if err != nil {
		return
	}

	logger.Info("\r\n\r\n%v say: ", pointToPointMessage.SourceUserName)
	logger.Notice("\t%v\n", pointToPointMessage.Content)
	return
}

func APPshowPointToPointMessage(responseMsg common.ResponseMessage, box *fyne.Container) (err error) {
	if responseMsg.Code != 200 {
		err = errors.New(responseMsg.Error)
		return
	}

	var pointToPointMessage common.PointToPointMessage
	err = json.Unmarshal([]byte(responseMsg.Data), &pointToPointMessage)
	if err != nil {
		return
	}

	logger.Info("\r\n\r\n%v say (APP): ", pointToPointMessage.SourceUserName)
	logger.Notice("\t%v\n", pointToPointMessage.Content)

	s := fmt.Sprintf("%v say: %v", pointToPointMessage.SourceUserName, pointToPointMessage.Content)
	text := widget.NewLabel(s)
	text.Wrapping = fyne.TextWrapBreak
	box.Add(text)

	return
}

// 处理服务端的返回
func Response(conn net.Conn, errMsg chan error) (err error) {
	var responseMsg common.ResponseMessage
	dispatcher := utils.Dispatcher{Conn: conn}

	for {
		responseMsg, err = dispatcher.ReadData()
		if err != nil {
			logger.Error("Waiting response error: %v\n", err)
			return
		}

		// 根据服务端返回的消息类型，进行相应的处理
		switch responseMsg.Type {
		case common.LoginResponseMessageType:
			err = dealLoginResponse(responseMsg)
			errMsg <- err
		case common.RegisterResponseMessageType:
			err = dealRegisterResponse(responseMsg)
			errMsg <- err
		case common.SendGroupMessageToClientType:
			err = dealGroupMessage(responseMsg)
			if err != nil {
				logger.Error("%v\n", err)
			}
		case common.ShowAllOnlineUsersType:
			err = showAllOnlineUsersList(responseMsg)
			errMsg <- err
		case common.PointToPointMessageType:
			err = showPointToPointMessage(responseMsg)
			errMsg <- err
		default:
			logger.Error("Unknown message type!")
		}

		if err != nil {
			return
		}
	}
}

// 处理服务端的返回
func APPResponse(conn net.Conn, errMsg chan error, userbox, groupbox, p2pbox *fyne.Container) (err error) {
	var responseMsg common.ResponseMessage
	dispatcher := utils.Dispatcher{Conn: conn}

	for {
		responseMsg, err = dispatcher.ReadData()
		if err != nil {
			logger.Error("Waiting response error: %v\n", err)
			return
		}

		// 根据服务端返回的消息类型，进行相应的处理
		switch responseMsg.Type {
		case common.LoginResponseMessageType:
			err = dealLoginResponse(responseMsg)
			errMsg <- err
		case common.RegisterResponseMessageType:
			err = dealRegisterResponse(responseMsg)
			errMsg <- err
		case common.SendGroupMessageToClientType:
			err = APPdealGroupMessage(responseMsg, groupbox)
			if err != nil {
				logger.Error("%v\n", err)
			}
		case common.ShowAllOnlineUsersType:
			err = APPshowAllOnlineUsersList(responseMsg, userbox)
			errMsg <- err
		case common.PointToPointMessageType:
			err = APPshowPointToPointMessage(responseMsg, p2pbox)
			errMsg <- err
		default:
			logger.Error("Unknown message type!")
		}

		if err != nil {
			return
		}
	}
}
