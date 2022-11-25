package model

import (
	"errors"
)

//根据业务逻辑需要，自定义一些错误
var (
	ErrorUserDoesNotExist = errors.New("user does not exist")
	ErrorUserPwd          = errors.New("password is invalid")

	// status code for register
	ErrorUserAlreadyExists    = errors.New("username already exists")
	ErrorPasswordDoesNotMatch = errors.New("password does not match")
)
