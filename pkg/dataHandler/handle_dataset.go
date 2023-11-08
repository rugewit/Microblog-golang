package dataHandler

import (
	"encoding/xml"
	"github.com/rugewit/microblog-golang/pkg/additinal"
	"github.com/rugewit/microblog-golang/pkg/models"
	"github.com/spf13/viper"
	"os"
)

func LoadUserAccounts() ([]*models.UserAccount, error) {
	defer additinal.Timer("LoadUserAccounts unmarshalling")()
	xmlFilePath := viper.Get("USERS_DATASET_PATH").(string)

	xmlData, err := os.ReadFile(xmlFilePath)
	if err != nil {
		return []*models.UserAccount{}, err
	}

	var users models.UserCollection

	err = xml.Unmarshal(xmlData, &users)
	if err != nil {
		return []*models.UserAccount{}, err
	}

	return users.Users.UserAccounts, nil
}

func LoadMessages() ([]*models.Message, error) {
	defer additinal.Timer("LoadMessages unmarshalling")()
	xmlFilePath := viper.Get("MESSAGES_DATASET_PATH").(string)

	xmlData, err := os.ReadFile(xmlFilePath)
	if err != nil {
		return []*models.Message{}, err
	}

	var messages models.MessageCollection

	err = xml.Unmarshal(xmlData, &messages)
	if err != nil {
		return []*models.Message{}, err
	}

	return messages.Posts.Messages, nil
}
