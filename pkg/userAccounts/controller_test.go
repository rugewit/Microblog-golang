package userAccounts

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"github.com/rugewit/microblog-golang/pkg/services"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func getDependenciesForUserService() (*services.MongoDataBase, context.Context) {
	envPath := "../env/.env"
	viper.SetConfigFile(envPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	dataBase, err := services.NewMongoDataBase()
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	return dataBase, ctx
}

func TestGetUserAccounts(t *testing.T) {
	dataBase, ctx := getDependenciesForUserService()
	userAccountService := services.NewUserAccountService(dataBase, ctx)
	router := gin.Default()
	RegisterRoutes(router, userAccountService)

	// request
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/users/", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
}
