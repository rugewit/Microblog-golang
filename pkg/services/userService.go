package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"github.com/rugewit/microblog-golang/pkg/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type UserAccountService struct {
	UserAccountCollection *mongo.Collection
	ctx                   context.Context
	rdb                   *redis.Client
	locker                *redislock.Client
}

func NewUserAccountService(database *MongoDataBase, ctx context.Context, collectionName string) *UserAccountService {
	userAccountService := new(UserAccountService)
	UserColName := viper.Get(collectionName).(string)
	var err error
	userAccountService.UserAccountCollection, err = database.GetMongoCollection(UserColName)
	if err != nil {
		panic(err)
	}
	userAccountService.ctx = ctx
	userAccountService.rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	userAccountService.rdb.FlushDB(userAccountService.ctx)
	userAccountService.locker = redislock.New(userAccountService.rdb)
	return userAccountService
}

func (service *UserAccountService) Create(userAccount *models.UserAccount) error {
	_, err := service.UserAccountCollection.InsertOne(service.ctx, userAccount)
	if err != nil {
		return err
	}
	return nil
}

func (service *UserAccountService) CreateMany(userAccounts []*models.UserAccount) error {
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

func (service *UserAccountService) Get(id string) (*models.UserAccount, error) {
	userAccount := new(models.UserAccount)

	// trying to find in redis
	redisResult, err := service.rdb.Get(service.ctx, id).Result()

	// not found in redis - load from mongo db
	if err == redis.Nil {
		log.Println("Load from mongo db")
		// load from mongo db
		objectId, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			return nil, err
		}

		res := service.UserAccountCollection.FindOne(service.ctx, bson.M{"_id": objectId})

		if res.Err() != nil {
			return nil, res.Err()
		}

		err = res.Decode(userAccount)
		if err != nil {
			return nil, err
		}

		// push into redis
		jsonRes, err := json.Marshal(userAccount)
		if err != nil {
			return nil, err
		}
		service.rdb.Set(service.ctx, id, jsonRes, 120*time.Second)
		return userAccount, nil

	} else if err != nil {
		return nil, err

	} else {
		log.Println("Load from redis")
		err := json.Unmarshal([]byte(redisResult), userAccount)
		if err != nil {
			return nil, err
		}
		return userAccount, nil
	}
}

func (service *UserAccountService) GetMany(limit int) ([]models.UserAccount, error) {
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

var IsLockedErr error = errors.New("UserAccount is locked")

func (service *UserAccountService) Update(id string, newUserAccount *models.UserAccount) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	lockTime := 60 * time.Second
	redactionSimulationTime := 30 * time.Second

	lockId := id + "1"
	errCh := make(chan error)
	go func() {
		lock, err := service.locker.Obtain(service.ctx, lockId, lockTime, nil)
		// it's already locked
		if errors.Is(err, redislock.ErrNotObtained) {
			fmt.Printf("userService %s is locked!\n", lockId)
			errCh <- IsLockedErr
			return
			// some other error
		} else if err != nil {
			errCh <- err
			return
		}
		errCh <- nil
		// Don't forget to defer Release.
		defer lock.Release(service.ctx)
		time.Sleep(redactionSimulationTime)
		_, err = service.UserAccountCollection.ReplaceOne(service.ctx, bson.M{"_id": objectId}, newUserAccount)
		service.rdb.Del(service.ctx, id)
	}()

	return <-errCh
}

func (service *UserAccountService) Delete(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	_, err = service.UserAccountCollection.DeleteOne(service.ctx, bson.M{"_id": objectId})
	if err != nil {
		return err
	}
	if isInCache := service.rdb.Get(service.ctx, id); isInCache != nil {
		service.rdb.Del(service.ctx, id)
	}
	return nil
}
