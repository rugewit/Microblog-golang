package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rugewit/microblog-golang/pkg/additinal"
	"github.com/rugewit/microblog-golang/pkg/dataHandler"
	"github.com/rugewit/microblog-golang/pkg/models"
	"github.com/rugewit/microblog-golang/pkg/services"
	"net/http"
)

type MessageController struct {
	messageService *services.MessageService
}

func NewMessageController(messageService *services.MessageService) *MessageController {
	return &MessageController{messageService: messageService}
}

func MessageRegisterRoutes(r *gin.Engine, messageService *services.MessageService) {
	messageController := NewMessageController(messageService)

	routes := r.Group("/messages")
	routes.POST("/", messageController.AddMessage)
	routes.GET("/:id", messageController.GetMessage)
	routes.GET("/", messageController.GetMessages)
	routes.GET("/load-messages", messageController.LoadMessages)
	routes.PUT("/:id", messageController.UpdateMessage)
	routes.DELETE("/:id", messageController.DeleteMessage)
}

func (messageController MessageController) AddMessage(c *gin.Context) {
	message := new(models.Message)

	// get request body
	if err := c.BindJSON(message); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := messageController.messageService.Create(message); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusCreated, message)
}

func (messageController MessageController) GetMessage(c *gin.Context) {
	id := c.Param("id")
	var message *models.Message
	var err error
	if message, err = messageController.messageService.Get(id); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, message)
}

func (messageController MessageController) GetMessages(c *gin.Context) {
	var messages []models.Message
	var err error

	limit := 200
	if messages, err = messageController.messageService.GetMany(limit); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}
	c.JSON(http.StatusOK, &messages)
}

func (messageController MessageController) LoadMessages(c *gin.Context) {
	var err error
	messages, err := dataHandler.LoadMessages()

	if err != nil {
		c.AbortWithError(http.StatusBadGateway, err)
	}

	defer additinal.Timer("LoadMessages pushing into mongo db")()

	if err = messageController.messageService.CreateMany(messages); err != nil {
		c.AbortWithError(http.StatusBadGateway, err)
	}

	c.JSON(http.StatusOK, "messages have been loaded")
}

func (messageController MessageController) UpdateMessage(c *gin.Context) {
	updatedMessage := new(models.Message)

	// get request body
	if err := c.BindJSON(updatedMessage); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id := c.Param("id")

	_, err := messageController.messageService.Get(id)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := messageController.messageService.Update(id, updatedMessage); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, updatedMessage)
}

func (messageController MessageController) DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	if err := messageController.messageService.Delete(id); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
