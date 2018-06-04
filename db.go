package main

import (
	"errors"
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	session     *mgo.Session
	heartColl   *mgo.Collection
	usersColl   *mgo.Collection
	userAccount = make(map[string]string)
)

//ID       bson.ObjectId `bson:"_id"`
type userDB struct {
	Username string `bson:"username"`
	Password []byte `bson:"password"`
	Slot     []byte `bson:"slot"`
	NickName string `bson:"nickname"`
	Email    string `bson:"email"`
}

type heartDB []struct {
	User  string
	Value int    `json:"value" binding:"required"`
	Time  uint64 `json:"time" binding:"required"`
}

func (h heartDB) toInterface(user string) (res []interface{}) {
	for _, v := range h {
		v.User = user
		res = append(res, v)
	}
	return res
}

func initDB() bool {
	//url := "mongodb://wdq:y1218@159.65.11.201:27017/love"
	dialInfo := &mgo.DialInfo{
		Addrs:     []string{"159.65.11.201"},
		Direct:    false,
		Timeout:   time.Second * 1,
		Database:  "love",
		Username:  "wdq",
		Password:  "y1218",
		PoolLimit: 4096,
	}
	var err error
	session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println("db err:", err)
		return false
	}
	db := session.DB("love")
	if nil == db {
		fmt.Println("db nil")
		return false
	}
	heartColl = db.C("heart")
	if nil == heartColl {
		fmt.Println("heartColl nil")
		return false
	}
	usersColl = db.C("users")
	if nil == usersColl {
		fmt.Println("usersColl nil")
		return false
	}

	return true
}

func insertUser(user userDB) (err error) {
	_, err = selectUserByUsername(user.Username)
	if nil != err {
		err = usersColl.Insert(&user)
		if nil != err {
			return
		} else {
			return
		}
	}
	return errors.New("other error")
}

func selectUserByUsername(username string) (user userDB, err error) {
	err = usersColl.Find(bson.M{"username": username}).One(&user)
	return
}

/*func getAllUsers() (users map[string]string) {
	users = make(map[string]string)
	one := userDB{}
	iter := usersColl.Find(nil).Iter()
	for iter.Next(one) {
		users[one.Username] = one.Password
	}
	return
}*/
