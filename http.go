package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	cookieMaxAge = time.Hour * 1 / time.Second
)

type userRequrest struct {
	Username string `bson:"username" binding:"required"`
	Password string `bson:"password"`
	NickName string `bson:"nickname"`
	Email    string `bson:"email"`
}

type userResponse struct {
	Username string `bson:"username"`
	Nickname string `bson:"nickname"`
}

type heartRequrest []struct {
	User  string
	Value int    `json:"value" binding:"required"`
	Time  uint64 `json:"time" binding:"required"`
}

func (h heartRequrest) toInterface(user string) (res []interface{}) {
	for _, v := range h {
		v.User = user
		res = append(res, v)
	}
	return res
}

func createUser(c *gin.Context) {
	var userReq userRequrest
	if nil == c.ShouldBind(&userReq) {
		slot, pa, err := genSlotAndHash(userReq.Password)
		if nil == err {
			var userdb userDB
			userdb.NickName = userReq.NickName
			userdb.Email = userReq.Email
			userdb.Username = userReq.Username
			userdb.Password = pa
			userdb.Slot = slot
			err = insertUser(userdb)
			if nil == err {
				c.JSON(http.StatusOK, gin.H{"status": "create success"})
				return
			} else {
				fmt.Println("err0:", err)
			}
		} else {
			fmt.Println("err1:", err)
		}
	} else {
		fmt.Println("err2:", "bad bind")
	}
	c.JSON(400, gin.H{"err": "create error"})
}

func getSlot(c *gin.Context) {
	c.JSON(200, gin.H{"slot": "get"})
}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		who, err := c.Cookie("who")
		alive, err := c.Cookie("logined")
		if nil != err {
			var userReq userRequrest
			if nil == c.ShouldBind(&userReq) {
				user, err := selectUserByUsername(userReq.Username)
				if nil == err {
					pa, err := calReqHash(user.Slot, userReq.Password)
					if nil == err && comper(pa, user.Password) {
						c.SetCookie("logined", "alive", int(cookieMaxAge), "/", "", false, true)
						c.SetCookie("who", userReq.Username, int(cookieMaxAge), "/", "", false, true)
						c.Next()
						return
					}
				}
			}
			c.AbortWithStatus(http.StatusUnauthorized)
		} else if "alive" == alive {
			c.SetCookie("logined", "alive", int(cookieMaxAge), "/", "", false, true)
			c.SetCookie("who", who, int(cookieMaxAge), "/", "", false, true)
			c.Next()
			return
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func push() gin.HandlerFunc {
	return func(c *gin.Context) {
		who, err := c.Cookie("who")
		if nil != err {
			c.Status(http.StatusFailedDependency)
			return
		}
		var heart heartRequrest
		if c.ShouldBind(&heart) == nil {
			/*if 0 == heart.Value {
				c.JSON(200, gin.H{"s": heart.Value})
			}*/
			err := heartColl.Insert(heart.toInterface(who))
			if err != nil {
				fmt.Println("insert ", err)
			}
			c.JSON(http.StatusOK, gin.H{"status": heart})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"ok": heart})
		}
	}
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})

	r.GET("/set", func(c *gin.Context) {
		c.SetCookie("tcoo", "yoyo", int(cookieMaxAge), "/", "", false, true)
		c.JSON(http.StatusOK, "set")
	})

	r.GET("/get", func(c *gin.Context) {
		tt, err := c.Cookie("tcoo")
		fmt.Println(tt, err)
		c.JSON(http.StatusOK, "get:"+tt)
	})

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		//fmt.Println(c, err)
		c.String(http.StatusOK, "pong")
	})

	// Create user
	r.POST("/user/create", createUser)

	// Get Slot
	r.POST("/user/slot", getSlot)

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		username := c.Params.ByName("name")
		user, err := selectUserByUsername(username)
		if nil == err {
			fmt.Println("user:", user)
			c.JSON(200, gin.H{"user": user.Username, "value": user})
		} else {
			c.JSON(200, gin.H{"user": user, "status": "no value"})
		}
	})

	//userAccount = getAllUsers()

	v1 := r.Group("/api/v1", auth())

	v1.POST("/login", func(c *gin.Context) {
		var userReq userRequrest
		err := c.ShouldBind(&userReq)
		if nil == err {
			user, err := selectUserByUsername(userReq.Username)
			if nil == err {
				var res userResponse
				res.Username = user.Username
				res.Nickname = user.NickName
				c.JSON(http.StatusOK, res)
			}
		} else {
			fmt.Println("login err:", err)
			c.Status(http.StatusBadRequest)
		}
	})

	v1.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, "test!")
	})

	//v1.POST("/push", push)

	//authorized := r.Group("/", gin.BasicAuth(userAccount))
	authorized := r.Group("/auth", auth())

	authorized.GET("gg", func(c *gin.Context) {
		fmt.Println("gg")
		c.JSON(200, gin.H{"ok": "okk"})
	})

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		var heart heartRequrest

		if c.ShouldBind(&heart) == nil {
			/*if 0 == heart.Value {
				c.JSON(200, gin.H{"s": heart.Value})
			}*/
			fmt.Println(user, "value:", heart)
			err := heartColl.Insert(heart.toInterface(user))
			if err != nil {
				fmt.Println("insert ", err)
			}
			c.JSON(200, gin.H{"status": heart})
		} else {
			c.JSON(200, gin.H{"ok": heart})
		}
	})

	return r
}
