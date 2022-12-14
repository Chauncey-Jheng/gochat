package process

import (
	"encoding/json"
	"fmt"
	common "go-chat/common/message"
	"go-chat/server/model"
	"go-chat/server/utils"
	"net"
)

type UserProcess struct {
	Conn net.Conn
}

// 响应客户端
func (this *UserProcess) responseClient(responseMessageType string, code int, data string) (err error) {
	var responseMessage common.ResponseMessage
	responseMessage.Code = code
	responseMessage.Type = responseMessageType
	responseMessage.Data = data

	responseData, err := json.Marshal(responseMessage)
	if err != nil {
		fmt.Printf("some error when generate response message, error: %v", err)
	}

	dispatcher := utils.Dispatcher{Conn: this.Conn}

	err = dispatcher.WriteData(responseData)
	return
}

func (this *UserProcess) UserRegister(message string) (err error) {
	var info common.RegisterMessage
	var code int
	data := ""
	err = json.Unmarshal([]byte(message), &info)
	if err != nil {
		code = common.ServerError
	}

	_, err = model.CurrentUserDao.Register(info.UserName, info.Password, info.PasswordConfirm)
	switch err {
	case nil:
		code = common.RegisterSucceed
	case model.ErrorPasswordDoesNotMatch:
		code = 402
	case model.ErrorUserAlreadyExists:
		code = 403
	default:
		code = 500
	}
	err = this.responseClient(common.RegisterResponseMessageType, code, data)
	return
}

func (this *UserProcess) UserLogin(message string) (err error) {
	var info common.LoginMessage
	var code int
	var data string
	err = json.Unmarshal([]byte(message), &info)
	if err != nil {
		code = common.ServerError
	}

	user, err := model.CurrentUserDao.Login(info.UserName, info.Password)

	switch err {
	case nil:
		code = common.LoginSucceed
		// save user conn status
		clientConn := model.ClientConn{}
		clientConn.Save(user.ID, user.Name, this.Conn)

		userInfo := common.UserInfo{ID: user.ID, UserName: user.Name}
		info, _ := json.Marshal(userInfo)
		data = string(info)
	case model.ErrorUserDoesNotExist:
		code = 404
	case model.ErrorUserPwd:
		code = 403
	default:
		code = 500
	}
	err = this.responseClient(common.LoginResponseMessageType, code, data)
	return
}
