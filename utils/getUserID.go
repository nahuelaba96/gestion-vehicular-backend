package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserID(c *gin.Context) (primitive.ObjectID, error) {
	val, exists := c.Get("user_id")
	if !exists {
		return primitive.NilObjectID, errors.New("no se encontró el user_id en el contexto")
	}

	userID, ok := val.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("formato de user_id inválido")
	}

	return userID, nil
}
