package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func createUser(c *gin.Context) {
	var user userStruct
	if nil == c.ShouldBind(&user) {
		err := insertUser(user)
		if nil != err {
			c.JSON(400, gin.H{"err": err.Error})
		} else {
			c.JSON(200, gin.H{"status": "success"})
		}
	} else {
		c.JSON(400, gin.H{"err": "contents error"})
	}
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Create user
	r.POST("/user/create", createUser)

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := DB[user]
		if ok {
			c.JSON(200, gin.H{"user": user, "value": value})
		} else {
			c.JSON(200, gin.H{"user": user, "status": "no value"})
		}
	})

	authorized := r.Group("/", gin.BasicAuth(userAccount))

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			fmt.Println("value:", json.Value)
			heartColl.Insert(&json)
			res := ""
			err := heartColl.Find(bson.M{}).One(&res)
			if err != nil {
				fmt.Println("find ", err)
			}
			fmt.Println("res:", res)
			DB[user] = json.Value
			c.JSON(200, gin.H{"status": json.Value})
		} else {
			c.JSON(200, gin.H{"ok": json.Value})
		}
	})

	return r
}
