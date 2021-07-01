package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/Nguyen-Hoang-Nam/let-shorten/db"
	"github.com/Nguyen-Hoang-Nam/let-shorten/form"
	"github.com/Nguyen-Hoang-Nam/let-shorten/models"
	"github.com/gin-gonic/gin"
)

var urlModel = new(models.Url)

func GetUrl(c *gin.Context) {
	var err error
	var url *models.Url

	url, err = urlModel.GetURLByHash(c.Param("id"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else if url != nil {
		c.String(http.StatusOK, url.URL)
	} else {
		c.String(http.StatusNotFound, "Url not found")
	}
}

func PostUrl(c *gin.Context) {
	var json form.Url
	var err error
	var cookie string

	client := db.GetRedis()

	cookie, err = c.Cookie("gin_cookie")
	if err != nil {
		c.String(http.StatusNotFound, "Not login user")
	} else {
		if err = c.ShouldBind(&json); err != nil {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			var id string
			id, err = client.Get(ctx, cookie).Result()
			if err != nil {
				c.String(http.StatusNotFound, err.Error())
			} else {
				hash := sha256.Sum256([]byte(json.URL))
				newUrl := models.Url{
					Hash: hex.EncodeToString(hash[:]),
					URL:  json.URL,
					TTL:  json.TTL,
				}

				if err = userModel.AddURLByID(newUrl, id); err != nil {
					c.String(http.StatusInternalServerError, err.Error())
				} else {
					newUrl := fmt.Sprintf("/url/%x", hash)
					c.String(http.StatusOK, newUrl)
				}
			}
		}
	}
}

func DeleteUrl(c *gin.Context) {
	var err error
	var cookie string

	position := c.Query("position")

	client := db.GetRedis()

	cookie, err = c.Cookie("gin_cookie")
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
	} else {
		var id string
		id, err = client.Get(ctx, cookie).Result()
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			if err = userModel.RemoveURLByID(id, position); err != nil {
				c.String(http.StatusInternalServerError, err.Error())
			} else {
				if err = urlModel.RemoveUrlByHash(c.Param("id")); err != nil {
					c.String(http.StatusInternalServerError, err.Error())
				} else {
					c.String(http.StatusOK, "Delete successful")
				}
			}
		}
	}
}
