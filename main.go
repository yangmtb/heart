package main

import (
	"fmt"
)

func main() {
	if !initDB() {
		fmt.Println("db error")
		return
	}
	//r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	//r.Run(":8080")
	session.Close()
}
