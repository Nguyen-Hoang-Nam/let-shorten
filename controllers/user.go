package controllers

import (
	"net/http"

	"github.com/Nguyen-Hoang-Nam/let-shorten/models"
	"github.com/gin-gonic/gin"
)

var userModel = new(models.User)

func GetUser(c *gin.Context) {
	user, err := userModel.GetByID(c.Param("id"))
	if err != nil {
		c.String(http.StatusNotFound, "Can not found user")
	} else {
		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}

func DeleteUser(c *gin.Context) {
	err := userModel.DeleteUser(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		c.String(http.StatusOK, "Delete successful")
	}
}
