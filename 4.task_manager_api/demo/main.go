package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type album struct {
	ID     string  `json:"id" validate:"required"`
	Title  string  `json:"title" validate:"required"`
	Artist string  `json:"artist" validate:"required"`
	Price  float64 `json:"price" validate:"required"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	fmt.Println("Hello, World!")
	router := gin.Default()
	router.GET("/test", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, "test")
	})
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, "ping")
	})

	// actual routes
	router.POST("/albums", createAlbum)
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumById)

	router.Run("localhost:8080")
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func createAlbum(c *gin.Context) {
	var newAlbum album
	// Call BindJSON to bind the received JSON to
	// newAlbum.
	// Initialize a validator
	v := validator.New()
	// data binding is too simple i dont' like it
	if err := c.BindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing album data"})
		// with body
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create album"})
		// no body
		// c.Status(http.StatusInternalServerError)
		return
	}
	// Validate the newAlbum struct
	if err := v.Struct(newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid album data"})
		return
	}

	fmt.Println(newAlbum)
	albums = append(albums, newAlbum)

	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")
	for _, val := range albums {
		if val.ID == id {
			c.IndentedJSON(http.StatusOK, val)
			return
		}
	}
	c.Status(http.StatusNoContent)
}
