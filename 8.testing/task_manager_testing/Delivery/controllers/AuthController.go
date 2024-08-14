package controllers

import (
	"fmt"

	"main/Domain"

	// "main/utils"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

type authController struct {
	AuthUseCase Domain.AuthUseCase
}
type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

func NewAuthController(service_reference Domain.AuthUseCase) (*authController, error) {

	// service_reference, err := UseCases.NewAuthUseCase(client, DataBase, _collection)
	// if err != nil {
	// 	return nil, err
	// }
	return &authController{
		AuthUseCase: service_reference,
	}, nil
}

func (ac *authController) Login(c *gin.Context) {
	var newUser Domain.User
	v := validator.New()
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid data", "error": err.Error()})
		return
	}
	if err := v.Struct(newUser); err != nil {
		fmt.Printf(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid or missing data", "error": err.Error()})
		return
	}
	token, err, statusCode := ac.AuthUseCase.Login(c, &newUser)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		//success
		c.IndentedJSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": token})
	}
}

func (ac *authController) Register(c *gin.Context) {
	var newUser Domain.User
	v := validator.New()
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid data", "error": err.Error()})
		return
	}
	if err := v.Struct(newUser); err != nil {
		fmt.Printf(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid or missing data", "error": err.Error()})
		return
	}
	createdUser, err, statusCode := ac.AuthUseCase.Register(c, &newUser)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		//success
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": createdUser})
	}
}
