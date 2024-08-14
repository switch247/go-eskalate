package Domain

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" bson:"_id,omitempty"`
	Email    string             `json:"email" validate:"required"`
	Password string             `json:"password,omitempty" validate:"required"`
	Is_Admin bool               `json:"is_admin,omitempty" default:"false"`
	Tasks    []Task             `json:"tasks,omitempty"`
}

// this could have been handled in a better way but i was too lazy to do it
type OmitedUser struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" bson:"_id,omitempty"`
	Email    string             `json:"email" validate:"required"`
	Password string             `json:"-"`
	Is_Admin bool               `json:"is_admin" default:"false"`
	Tasks    []Task             `json:"tasks,omitempty"`
}

type AuthRepository interface {
	Login(ctx context.Context, user *User) (string, error, int)
	Register(ctx context.Context, user *User) (OmitedUser, error, int)
}

type AuthUseCase interface {
	Login(c *gin.Context, user *User) (string, error, int)
	Register(c *gin.Context, user *User) (OmitedUser, error, int)
}

type UserRepository interface {
	CreateUsers(ctx context.Context, user *User) (OmitedUser, error, int)
	GetUsers(ctx context.Context) ([]*OmitedUser, error, int)
	GetUsersById(ctx context.Context, id primitive.ObjectID, user OmitedUser) (OmitedUser, error, int)
	UpdateUsersById(ctx context.Context, id primitive.ObjectID, user User, curentuser OmitedUser) (OmitedUser, error, int)
	DeleteUsersById(ctx context.Context, id primitive.ObjectID, user OmitedUser) (error, int)
}

type UserUseCases interface {
	CreateUsers(c *gin.Context, user *User) (OmitedUser, error, int)
	GetUsers(c *gin.Context) ([]*OmitedUser, error, int)
	GetUsersById(c *gin.Context, id primitive.ObjectID, user OmitedUser) (OmitedUser, error, int)
	UpdateUsersById(c *gin.Context, id primitive.ObjectID, user User, curentuser OmitedUser) (OmitedUser, error, int)
	DeleteUsersById(c *gin.Context, id primitive.ObjectID, user OmitedUser) (error, int)
}
