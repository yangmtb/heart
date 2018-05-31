package main

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	session     *mgo.Session
	heartColl   *mgo.Collection
	usersColl   *mgo.Collection
	DB          = make(map[string]string)
	userAccount = make(map[string]string)
)

type userStruct struct {
	//ID       bson.ObjectId `bson:"_id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	NickName string `bson:"nickname"`
	email    string `bson:"email"`
}

type heartStruct struct {
	tim time.Time `bson:"time"`
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

func insertUser(user userStruct) (err error) {
	_, err = selectUserByUsername(user.Username)
	if nil != err {
		return
	}
	err = usersColl.Insert(&user)
	if nil != err {
		return
	}
	return
}

func selectUserByUsername(username string) (user userStruct, err error) {
	err = usersColl.Find(bson.M{"username": username}).One(&user)
	return
}

func getAllUsers() (users map[string]string) {
	users = make(map[string]string)
	one := userStruct{}
	iter := usersColl.Find(nil).Iter()
	for iter.Next(one) {
		users[one.Username] = one.Password
	}
	return
}
