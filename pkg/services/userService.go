package services

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/rugewit/microblog-golang/pkg/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type UserMongoService struct {
	UserAccountCollection *mongo.Collection
	ctx                   context.Context
	rdb                   *redis.Client
}

func NewUserAccountService(database *MongoDataBase, ctx context.Context) *UserMongoService {
	userMongoService := new(UserMongoService)
	UserColName := viper.Get("USER_COL_NAME").(string)
	userMongoService.UserAccountCollection = database.GetMongoCollection(UserColName)
	userMongoService.ctx = ctx
	userMongoService.rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return userMongoService
}

func (service *UserMongoService) Create(userAccount *models.UserAccount) error {
	_, err := service.UserAccountCollection.InsertOne(service.ctx, userAccount)
	if err != nil {
		return err
	}
	return nil
}

func (service *UserMongoService) CreateMany(userAccounts []*models.UserAccount) error {
	interfaceSlice := make([]interface{}, len(userAccounts))
	for i, account := range userAccounts {
		interfaceSlice[i] = account
	}

	_, err := service.UserAccountCollection.InsertMany(service.ctx, interfaceSlice)
	if err != nil {
		return err
	}
	return nil
}

func (service *UserMongoService) Get(id string) (*models.UserAccount, error) {
	userAccount := new(models.UserAccount)

	objectId, _ := primitive.ObjectIDFromHex(id)
	res := service.UserAccountCollection.FindOne(service.ctx, bson.M{"_id": objectId})

	if res.Err() != nil {
		return userAccount, res.Err()
	}

	err := res.Decode(userAccount)
	if err != nil {
		return userAccount, err
	}

	return userAccount, nil
}

func (service *UserMongoService) GetMany(limit int) ([]models.UserAccount, error) {
	var userAccounts []models.UserAccount

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))

	cur, err := service.UserAccountCollection.Find(service.ctx, bson.D{{}}, findOptions)

	if err != nil {
		return nil, err
	}

	for cur.Next(service.ctx) {
		userAccount := models.UserAccount{}
		if err := cur.Decode(&userAccount); err != nil {
			log.Fatal(err)
			return nil, err
		}
		userAccounts = append(userAccounts, userAccount)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = cur.Close(service.ctx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return userAccounts, nil
}

func (service *UserMongoService) Update(id string, newUserAccount *models.UserAccount) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = service.UserAccountCollection.ReplaceOne(service.ctx, bson.M{"_id": objectId}, newUserAccount)

	if err != nil {
		return err
	}

	return nil
}

func (service *UserMongoService) Delete(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	_, err = service.UserAccountCollection.DeleteOne(service.ctx, bson.M{"_id": objectId})
	if err != nil {
		return err
	}
	return nil
}
