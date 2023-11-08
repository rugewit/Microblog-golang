package services

import (
	"context"
	"errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoDataBase struct {
	Name     string
	Address  string
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDataBase() (*MongoDataBase, error) {
	db := new(MongoDataBase)
	ctx := context.TODO()

	dbName := viper.Get("DB_NAME").(string)
	dbAddress := viper.Get("DB_URL").(string)

	db.Name = dbName
	db.Address = dbAddress

	clientOptions := options.Client().ApplyURI(dbAddress)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	db.Client = client
	db.Database = client.Database(db.Name)
	return db, nil
}

func (db *MongoDataBase) GetMongoCollection(collectionName string) (*mongo.Collection, error) {
	collection := db.Database.Collection(collectionName)
	if collection == nil {
		return nil, errors.New("collection is nil")
	}
	return db.Database.Collection(collectionName), nil
}
