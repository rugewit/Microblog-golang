package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rugewit/microblog-golang/pkg/additinal"
	"github.com/rugewit/microblog-golang/pkg/dataHandler"
	"github.com/rugewit/microblog-golang/pkg/models"
	"github.com/rugewit/microblog-golang/pkg/services"
	"net/http"
)

type UserController struct {
	userAccountService *services.UserMongoService
}

func NewUserController(userAccountService *services.UserMongoService) *UserController {
	return &UserController{userAccountService: userAccountService}
}

func UserRegisterRoutes(r *gin.Engine, userAccService *services.UserMongoService) {
	userController := NewUserController(userAccService)

	routes := r.Group("/users")
	routes.POST("/", userController.AddUserAccount)
	routes.GET("/:id", userController.GetUserAccount)
	routes.GET("/", userController.GetUserAccounts)
	routes.GET("/load-users", userController.LoadUserAccounts)
	routes.PUT("/:id", userController.UpdateUserAccount)
	routes.DELETE("/:id", userController.DeleteUserAccounts)
}

func (userController UserController) AddUserAccount(c *gin.Context) {
	userAccount := new(models.UserAccount)

	// get request body
	if err := c.BindJSON(userAccount); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := userController.userAccountService.Create(userAccount); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusCreated, userAccount)
}

func (userController UserController) GetUserAccount(c *gin.Context) {
	id := c.Param("id")
	var userAccount *models.UserAccount
	var err error
	if userAccount, err = userController.userAccountService.Get(id); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, userAccount)
}

func (userController UserController) GetUserAccounts(c *gin.Context) {
	var userAccounts []models.UserAccount
	var err error

	limit := 200
	if userAccounts, err = userController.userAccountService.GetMany(limit); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}
	c.JSON(http.StatusOK, &userAccounts)
}

func (userController UserController) LoadUserAccounts(c *gin.Context) {
	var err error
	userAccounts, err := dataHandler.LoadUserAccounts()

	if err != nil {
		c.AbortWithError(http.StatusBadGateway, err)
	}

	defer additinal.Timer("LoadUserAccounts pushing into mongo db")()

	if err = userController.userAccountService.CreateMany(userAccounts); err != nil {
		c.AbortWithError(http.StatusBadGateway, err)
	}

	c.JSON(http.StatusOK, "users have been loaded")
}

func (userController UserController) UpdateUserAccount(c *gin.Context) {
	updatedUserAccount := new(models.UserAccount)

	// get request body
	if err := c.BindJSON(updatedUserAccount); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id := c.Param("id")

	_, err := userController.userAccountService.Get(id)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := userController.userAccountService.Update(id, updatedUserAccount); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, updatedUserAccount)
}

func (userController UserController) DeleteUserAccounts(c *gin.Context) {
	id := c.Param("id")
	if err := userController.userAccountService.Delete(id); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
