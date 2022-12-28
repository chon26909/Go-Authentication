package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection interface {
	Db() *mongo.Database
}

type conn struct {
	database *mongo.Database
}

// func getDbURL() string {
// 	port, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
// 	if err != nil {
// 		log.Println("Error on load db port from env :", err.Error())
// 		port = 27017
// 	}

// 	DATABASE_USER := os.Getenv("DATABASE_USER")
// 	DATABASE_PASSWORD := os.Getenv("DATABASE_PASS")
// 	DATABASE_HOST := os.Getenv("DATABASE_HOST")
// 	DATABASE_NAME := os.Getenv("DATABASE_NAME")

// 	return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", DATABASE_USER, DATABASE_PASSWORD, DATABASE_HOST, port, DATABASE_NAME)
// }

func NewConnection() Connection {
	var c conn
	var uri string = os.Getenv("MONGODB_URI")

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	c.database = client.Database("DEV")

	return &c
}

func (c *conn) Db() *mongo.Database {
	return c.database
}
