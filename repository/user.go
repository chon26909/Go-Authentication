package repository

import (
	"context"
	"go-auth/db"
	"go-auth/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const UsersCollection = "users"

type UsersRepository interface {
	Save(user *models.User) error
	Update(user *models.User) error
	GetById(id string) (user *models.User, err error)
	GetByEmail(email string) (user *models.User, err error)
	GetAll() (users []*models.User, err error)
	Delete(id string) error
}

type usersRepository struct {
	c *mongo.Collection
}

func NewUserRepository(conn db.Connection) UsersRepository {
	return &usersRepository{conn.Db().Collection(UsersCollection)}
}

func (r *usersRepository) Save(user *models.User) error {
	_, err := r.c.InsertOne(context.TODO(), user)
	return err
}

func (r *usersRepository) Update(user *models.User) error {
	_, err := r.c.UpdateByID(context.TODO(), user.Id, user)
	return err
}

func (r *usersRepository) GetById(uid string) (user *models.User, err error) {

	var result *models.User

	objectId, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objectId}}
	err = r.c.FindOne(context.TODO(), filter).Decode(&result)

	return result, err
}

func (r *usersRepository) GetByEmail(email string) (user *models.User, err error) {

	var result *models.User

	filter := bson.D{{Key: "email", Value: email}}
	err = r.c.FindOne(context.TODO(), filter).Decode(&result)

	return result, err
}

func (r *usersRepository) GetAll() (users []*models.User, err error) {

	filter := bson.D{{}}
	cursor, err := r.c.Find(context.TODO(), filter)

	if err != nil {
		return nil, err
	}

	var resultUsers []*models.User
	cursor.All(context.TODO(), &resultUsers)

	return resultUsers, err
}

func (r *usersRepository) Delete(id string) error {

	filter := bson.D{{Key: "_id", Value: id}}

	_, err := r.c.DeleteOne(context.TODO(), filter)
	return err
}
