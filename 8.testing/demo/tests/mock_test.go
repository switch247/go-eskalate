package tests

import (
	"fmt"

	"testing"

	"github.com/stretchr/testify/mock"
)

// Mocking

// ShapeServiceMock mocks the ShapeService interface
type ShapeServiceMock struct {
	mock.Mock
}

func (m *ShapeServiceMock) CalculateArea(radius float64) float64 {
	fmt.Println("Mocked area calculation function")
	fmt.Printf("Radius passed in: %f\n", radius)
	args := m.Called(radius)
	return args.Get(0).(float64)
}

func (m *ShapeServiceMock) DummyFunc() {
	fmt.Println("Dummy")
}

// CircleService represents a service for circle-related calculations
type CircleService struct {
	shapeService *ShapeServiceMock
}

// CalculateCircleArea calculates the area of a circle using the provided radius
func (cs CircleService) CalculateCircleArea(radius float64) float64 {
	return cs.shapeService.CalculateArea(radius)
}

func TestCalculateCircleArea(t *testing.T) {
	shapeMock := new(ShapeServiceMock)
	expectedArea := 78.54
	shapeMock.On("CalculateArea", 5.0).Return(expectedArea)

	circleService := CircleService{shapeService: shapeMock}
	result := circleService.CalculateCircleArea(5.0)

	// Verify that the expectations were met
	shapeMock.AssertExpectations(t)

	// Additional assertion for the calculated area
	if result != expectedArea {
		t.Errorf("Expected area %f, but got %f", expectedArea, result)
	}
}

func TestDummyFunc(t *testing.T) {
	shapeMock := new(ShapeServiceMock)
	shapeMock.DummyFunc()
	// Add assertions or verifications for the behavior of DummyFunc if needed
}
