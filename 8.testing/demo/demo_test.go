package demo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mock *MockRepository) Create(i int) (int, error) {
	log.Fatal("is this even being called?")
	args := mock.Called(i)
	return args.Int(0), args.Error(1)
}

func TestCreate(t *testing.T) {
	mockRepo := new(MockRepository)

	mockRepo.On("Create", 1).Return(1, nil)

	demoservice := NewRepository(mockRepo)

	result, err := demoservice.Create(1)

	mockRepo.AssertExpectations(t)

	assert.Equal(t, 1, result)
	assert.Nil(t, err)
}
