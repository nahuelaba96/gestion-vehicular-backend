package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Email    string             `bson:"email"`
    Name     string             `bson:"name"`
    OAuthID  string             `bson:"oauth_id"` // opcional
}
