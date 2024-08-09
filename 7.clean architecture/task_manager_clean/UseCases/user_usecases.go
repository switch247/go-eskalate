package UseCases

import (
	"context"
	"main/Domain"
	"main/Repositories"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userUseCases struct {
	userRepository Domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUseCase(client *mongo.Client, DataBase *mongo.Database, _collection *mongo.Collection) (*userUseCases, error) {
	service_reference, err := Repositories.NewUserRepository(client, DataBase, _collection)
	if err != nil {
		return nil, err
	}
	return &userUseCases{
		userRepository: service_reference,
		contextTimeout: time.Second * 10,
	}, nil
}

func (uc *userUseCases) GetUsers(c *gin.Context) ([]*Domain.OmitedUser, error, int) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.userRepository.GetUsers(ctx)

}

func (uc *userUseCases) GetUsersById(c *gin.Context, id primitive.ObjectID, user Domain.OmitedUser) (Domain.OmitedUser, error, int) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.userRepository.GetUsersById(ctx, id, user)

}

func (uc *userUseCases) CreateUsers(c *gin.Context, user *Domain.User) (Domain.OmitedUser, error, int) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.userRepository.CreateUsers(ctx, user)

}

func (uc *userUseCases) UpdateUsersById(c *gin.Context, id primitive.ObjectID, user Domain.User, curentuser Domain.OmitedUser) (Domain.OmitedUser, error, int) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.userRepository.UpdateUsersById(ctx, id, user, curentuser)

}

func (uc *userUseCases) DeleteUsersById(c *gin.Context, id primitive.ObjectID, user Domain.OmitedUser) (error, int) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.userRepository.DeleteUsersById(ctx, id, user)

}
