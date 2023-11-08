package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rugewit/microblog-golang/pkg/controllers"
	"github.com/rugewit/microblog-golang/pkg/services"
	"github.com/spf13/viper"
	"log"
)

func LoadViper(envPath string) error {
	viper.SetConfigFile(envPath)
	err := viper.ReadInConfig()

	if err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Println("Hello backend!")

	err := LoadViper("./pkg/env/.env")
	if err != nil {
		log.Fatal(err)
		return
	}

	port := viper.Get("PORT").(string)

	if err != nil {
		log.Fatal(err)
		return
	}

	dataBase, err := services.NewMongoDataBase()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	userAccountService := services.NewUserAccountService(dataBase, ctx)
	messageService := services.NewMessageService(dataBase, ctx)

	router := gin.Default()
	controllers.UserRegisterRoutes(router, userAccountService)
	controllers.MessageRegisterRoutes(router, messageService)
	router.Run(port)
}
