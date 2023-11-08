package services

import (
	"context"
	"github.com/rugewit/microblog-golang/pkg/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageService struct {
	MessageCollection *mongo.Collection
	ctx               context.Context
}

func NewMessageService(database *MongoDataBase, ctx context.Context) *MessageService {
	messageService := new(MessageService)
	UserColName := viper.Get("MESSAGE_COL_NAME").(string)
	var err error
	messageService.MessageCollection, err = database.GetMongoCollection(UserColName)
	if err != nil {
		panic(err)
	}
	messageService.ctx = ctx
	return messageService
}

func (service *MessageService) Create(message *models.Message) error {
	_, err := service.MessageCollection.InsertOne(service.ctx, message)
	if err != nil {
		return err
	}
	return nil
}

func (service *MessageService) CreateMany(messages []*models.Message) error {
	interfaceSlice := make([]interface{}, len(messages))
	for i, account := range messages {
		interfaceSlice[i] = account
	}

	_, err := service.MessageCollection.InsertMany(service.ctx, interfaceSlice)
	if err != nil {
		return err
	}
	return nil
}

func (service *MessageService) Get(id string) (*models.Message, error) {
	// load from mongo db
	message := new(models.Message)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res := service.MessageCollection.FindOne(service.ctx, bson.M{"_id": objectId})

	if res.Err() != nil {
		return nil, res.Err()
	}

	err = res.Decode(message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (service *MessageService) GetMany(limit int) ([]models.Message, error) {
	var messages []models.Message

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))

	cur, err := service.MessageCollection.Find(service.ctx, bson.D{{}}, findOptions)

	if err != nil {
		return nil, err
	}

	for cur.Next(service.ctx) {
		message := models.Message{}
		if err := cur.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	err = cur.Close(service.ctx)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (service *MessageService) Update(id string, newMessage *models.Message) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = service.MessageCollection.ReplaceOne(service.ctx, bson.M{"_id": objectId}, newMessage)

	if err != nil {
		return err
	}

	return nil
}

func (service *MessageService) Delete(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	_, err = service.MessageCollection.DeleteOne(service.ctx, bson.M{"_id": objectId})
	if err != nil {
		return err
	}
	return nil
}
