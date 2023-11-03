package userAccounts

import (
	"github.com/gin-gonic/gin"
	"github.com/rugewit/microblog-golang/pkg/models"
	"github.com/rugewit/microblog-golang/pkg/services"
	"net/http"
)

type Handler struct {
	userAccountService *services.UserMongoService
}

func RegisterRoutes(r *gin.Engine, userAccService *services.UserMongoService) {
	h := &Handler{userAccService}

	routes := r.Group("/users")
	routes.POST("/", h.AddUserAccount)
	routes.GET("/:id", h.GetUserAccount)
	routes.GET("/", h.GetUserAccounts)
	routes.PUT("/:id", h.UpdateUserAccount)
	routes.DELETE("/:id", h.DeleteUserAccounts)
}

func (h Handler) AddUserAccount(c *gin.Context) {
	userAccount := new(models.UserAccount)

	// get request body
	if err := c.BindJSON(userAccount); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.userAccountService.Create(userAccount); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusCreated, userAccount)
}

func (h Handler) GetUserAccount(c *gin.Context) {
	id := c.Param("id")
	var userAccount *models.UserAccount
	var err error
	if userAccount, err = h.userAccountService.Get(id); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, userAccount)
}

func (h Handler) GetUserAccounts(c *gin.Context) {
	var userAccounts []models.UserAccount
	var err error

	limit := 200
	if userAccounts, err = h.userAccountService.GetMany(limit); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}
	c.JSON(http.StatusOK, &userAccounts)
}

func (h Handler) UpdateUserAccount(c *gin.Context) {
	updatedUserAccount := new(models.UserAccount)

	// get request body
	if err := c.BindJSON(updatedUserAccount); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id := c.Param("id")

	_, err := h.userAccountService.Get(id)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.userAccountService.Update(id, updatedUserAccount); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, updatedUserAccount)
}

func (h Handler) DeleteUserAccounts(c *gin.Context) {
	id := c.Param("id")
	if err := h.userAccountService.Delete(id); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
