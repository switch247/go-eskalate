package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/models"

	"github.com/gin-gonic/gin"
)

func ReadJSON(ctx *gin.Context) (map[string]interface{}, error) {
	var jsonData map[string]interface{}
	err := json.NewDecoder(ctx.Request.Body).Decode(&jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func ValidateStatus(t *models.Task) error {
	validStatus := []models.TaskStatus{models.StatusPending, models.StatusCompleted, models.StatusInProgress}
	for _, status := range validStatus {
		if t.Status == status {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Invalid status: %s", t.Status))
}
