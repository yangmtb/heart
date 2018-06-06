package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type userInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

func struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	fmt.Println(t.NumMethod())
	fmt.Println(v)
	data := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Field(i).Tag.Get("json"))
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func (obj *userInfo) toMap() (res map[string]interface{}) {
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
