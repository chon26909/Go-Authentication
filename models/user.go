package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
	Salt     string             `bson:"salt"`
	CreateAt time.Time          `bson:"create_at"`
	UpdateAt time.Time          `bson:"update_at"`
}
