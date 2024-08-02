package main

import (
	"fmt"
	"log"
	"main/router"
	// "main/utils"
	// "net/http"
	// "github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Task Manager API")
	r := router.NewRouter()
	if r != nil {
		r.Run()
	} else {
		log.Fatal("Failed to start server")
	}

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
