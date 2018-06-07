package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

var (
	userClt  *redis.Client
	heartClt *redis.Client
)

type UserInfo struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
	Slot     []byte `json:"slot"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

func (obj *UserInfo) toMap() (res map[string]interface{}) {
	e, err := json.Marshal(obj)
	if nil != err {
		fmt.Println("marshal err", err)
		return
	}
	err = json.Unmarshal(e, &res)
	if nil != err {
		fmt.Println("unmarshal err", err)
		return
	}
	return
}

func initDB() bool {
	userClt = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "wdq",
		DB:       0,
	})
	heartClt = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "wdq",
		DB:       1,
	})
	return true
}

func selectUserByUsername(username string) (user UserInfo, err error) {
	userMap, err := userClt.HGetAll(username).Result()
	if nil != err {
		return
	}
	//fmt.Println("map", userMap)
	userJSON, err := json.Marshal(userMap)
	if nil != err {
		return
	}
	//fmt.Println("json:", string(userJSON))
	err = json.Unmarshal(userJSON, &user)
	return
}

func insertUser(user UserInfo) (err error) {
	ex := userClt.Exists(user.Username)
	if 0 != ex.Val() {
		fmt.Println("ex:", ex.Val())
		return fmt.Errorf("username is exsit")
	}
	cmd := userClt.HMSet(user.Username, user.toMap())
	if "OK" != cmd.Val() {
		fmt.Println("cmd", cmd.Val())
		return fmt.Errorf("set err")
	}
	return
}
