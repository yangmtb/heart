package main

import (
	"fmt"
)

type ttDB struct {
	Test []byte `bson:"username"`
}

func main() {
	/*data := "password"
	senc := base64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(senc)
	sdec, err := base64.StdEncoding.DecodeString(senc)
	fmt.Println(err)
	fmt.Println(string(sdec))
	//return
	x, err := randByte(256)
	fmt.Println("y")
	y, err := randByte(256)
	if nil != err {
		fmt.Println(err)
	}
	fmt.Println("res:", comper(x, y))
	fmt.Println(x)
	fmt.Println(y)
	//return*/
	if !initDB() {
		fmt.Println("db error")
		return
	}
	defer session.Close()
	var tt ttDB
	var err error
	tt.Test, err = randByte(32)
	fmt.Println(err, tt)
	//err = usersColl.Insert(&tt)
	//fmt.Println(err)
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
