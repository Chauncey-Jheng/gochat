package model

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// UserDao (User Data access object)实例，全局唯一
var CurrentUserDao *UserDao

type UserDao struct {
	pool *redis.Pool
}

// 初始化一个 UserDao 结构体实例
func InitUserDao(pool *redis.Pool) (currentUserDao *UserDao) {
	currentUserDao = &UserDao{pool: pool}
	return
}

// 用户 id 增长
func idIncr(conn redis.Conn) (id int, err error) {
	res, err := conn.Do("incr", "users_id")
	if err != nil {
		fmt.Printf("id incr error: %v\n", err)
		return
	}
	id = int(res.(int64))
	return
}

// 根据用户 id 获取用户信息
// 获取成功返回 user 信息，err nil
// 获取失败返回 err，user 为 nil
func (this *UserDao) GetUsrById(id int) (user User, err error) {
	conn := this.pool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("hget", "users", id))
	if err != nil {
		err = ErrorUserDoesNotExist
		return
	}
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Printf("Unmarshal user info error: %v\n", err)
		return
	}
	return
}

// 根据用户 username 获取用户信息
// 获取成功返回 user 信息，err nil
// 获取失败返回 err，user 为 nil
func (this *UserDao) GetUserByUserName(userName string) (user User, err error) {
	conn := this.pool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("hget", "users", userName))
	if err != nil {
		err = ErrorUserDoesNotExist
		return
	}
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Printf("Unmarshal user info error: %v\n", err)
		return
	}
	return
}

// 注册用户
// 用户名不能重复
func (this *UserDao) Register(userName, password, passwordConfirm string) (user User, err error) {
	// 判断密码是否正确
	if password != passwordConfirm {
		err = ErrorPasswordDoesNotMatch
		return
	}

	// 保证用户名不重复
	user, err = this.GetUserByUserName(userName)
	if err == nil {
		fmt.Printf("User already exists!\n")
		err = ErrorUserAlreadyExists
		return
	}

	conn := this.pool.Get()
	defer conn.Close()
	// id 自增 1，作为下个用户 id
	id, err := idIncr(conn)
	if err != nil {
		return
	}

	user = User{ID: id, Name: userName, Password: password}
	info, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("json marshal error: %v", err)
		return
	}
	_, err = conn.Do("hset", "users", userName, info)
	if err != nil {
		fmt.Printf("set user to reids error: %v", err)
		return
	}
	return
}

func (this *UserDao) Login(userName, password string) (user User, err error) {
	user, err = this.GetUserByUserName(userName)
	if err != nil {
		fmt.Printf("get user by id error: %v\n", err)
		return
	}

	if user.Password != password {
		err = ErrorUserPwd
		return
	}

	return
}
