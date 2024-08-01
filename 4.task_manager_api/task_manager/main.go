package main

import (
	"fmt"
	"main/router"
	// "main/utils"
	// "net/http"
	// "github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Task Manager API")
	r := router.NewRouter()
	// r.GET("/ping", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	// r.POST("/test", testAnyjson)
	r.Run()

}

// func testAnyjson(ctx *gin.Context) {
// 	val, err := utils.ReadJSON(ctx)
// 	if err != nil {
// 		fmt.Println("failed")
// 		fmt.Println(err.Error())
// 		ctx.Status(http.StatusBadRequest)

// 	} else {
// 		fmt.Println(val["id"])
// 		if val["xid"] != nil {
// 			ctx.JSON(200, val["xid"])
// 		}
// 		return
// 	}
// }
