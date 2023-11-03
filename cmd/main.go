package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rugewit/microblog-golang/pkg/models"
	"github.com/rugewit/microblog-golang/pkg/services"
	"github.com/rugewit/microblog-golang/pkg/userAccounts"
	"github.com/spf13/viper"
	"log"
)

func main() {
	fmt.Println("Hello backend!")
	envPath := "./pkg/env/.env"
	viper.SetConfigFile(envPath)
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
		return
	}
	port := viper.Get("PORT").(string)

	dataBase, err := services.NewMongoDataBase()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	userAccountService := services.NewUserAccountService(dataBase, ctx)

	_ = models.UserAccount{
		Reputation:     4,
		CreationDate:   "123",
		DisplayName:    "Alexey",
		LastAccessDate: "yesterday",
		WebsiteUrl:     "ya.ru",
		Location:       "Moscow",
		AboutMe:        "I am not superman",
		Views:          12,
		UpVotes:        13,
		DownVotes:      5,
		AccountId:      84,
	}

	router := gin.Default()
	userAccounts.RegisterRoutes(router, userAccountService)
	router.Run(port)
}
