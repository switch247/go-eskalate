package UseCases

import (
	"context"
	"time"

	"main/Domain"
	"main/Repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type authUseCase struct {
	AuthRepository Domain.AuthRepository
	contextTimeout time.Duration
}

func NewAuthUseCase(client *mongo.Client, DataBase *mongo.Database, _collection *mongo.Collection) (Domain.AuthUseCase, error) {
	service_reference, err := Repositories.NewAuthRepository(client, DataBase, _collection)
	if err != nil {
		return nil, err
	}
	return &authUseCase{
		AuthRepository: service_reference,
		contextTimeout: time.Second * 10,
	}, nil
}

func (au *authUseCase) Login(c *gin.Context, user *Domain.User) (string, error, int) {
	ctx, cancel := context.WithTimeout(c, au.contextTimeout)
	defer cancel()

	return au.AuthRepository.Login(ctx, user)

}

func (au *authUseCase) Register(c *gin.Context, user *Domain.User) (Domain.OmitedUser, error, int) {
	ctx, cancel := context.WithTimeout(c, au.contextTimeout)
	defer cancel()
	return au.AuthRepository.Register(ctx, user)

}
