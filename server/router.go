package server

import (
	"net/http"

	"github.com/Nguyen-Hoang-Nam/let-shorten/controllers"
	"github.com/Nguyen-Hoang-Nam/let-shorten/models"
	"github.com/gin-gonic/gin"
)

var urlModel = new(models.Url)

func NewRouter() *gin.Engine {
	// Middleware
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/s/:id", controllers.GetUrl)

	router.POST("/login", controllers.Login)

	router.POST("/signup", controllers.Signup)

	router.POST("/forget-password", func(c *gin.Context) {
		c.String(http.StatusOK, "Forget password")
	})

	router.POST("/reset-password/:resetid", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"id": c.Param("resetid"),
		})
	})

	router.POST("/new-password/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"id": c.Param("id"),
		})
	})

	router.GET("/signout", controllers.Signout)

	router.GET("/user/:id", controllers.GetUser)

	router.DELETE("/user/:id", controllers.DeleteUser)

	router.GET("/url/:id", controllers.GetUrl)

	router.POST("/url", controllers.PostUrl)

	router.PUT("/url/:urlid", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"url": c.Param("urlid"),
		})
	})

	router.DELETE("/url/:urlid", controllers.DeleteUrl)

	return router
}
