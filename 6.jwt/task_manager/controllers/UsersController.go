package controllers

import (
	// "fmt"
	// "errors"
	// "main/models"
	"main/data"
	"main/models"
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController interface {
	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userController struct {
	authService data.UserService
}

func NewUserController() (*userController, error) {
	service_reference, err := data.NewUserService()
	if err != nil {
		return nil, err
	}
	return &userController{
		authService: *service_reference,
	}, nil
}
func (uc *userController) GetUsers(c *gin.Context) {

	users, err, statusCode := uc.authService.GetUsers()
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(statusCode, gin.H{"users": users})
	}
}

func (uc *userController) GetUser(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	user, err, statusCode := uc.authService.GetUsersById(objectID, user)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(statusCode, gin.H{"user": user})
	}
}

func (uc *userController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err, statusCode := uc.authService.CreateUsers(&user)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(statusCode, gin.H{"user": user})
	}
}

func (uc *userController) UpdateUser(c *gin.Context) {
	logeduser, err := utils.ExtractUser(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(400, gin.H{"error": err.Error()})
		return
	}
	OmitedUser, err, statusCode := uc.authService.UpdateUsersById(objectID, user, logeduser)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(statusCode, gin.H{"user": OmitedUser})
	}
}

func (uc *userController) DeleteUser(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	err, statusCode := uc.authService.DeleteUsersById(objectID, user)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(statusCode, gin.H{"message": "User deleted successfully"})
	}
}
