package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Nguyen-Hoang-Nam/let-shorten/db"
	"github.com/Nguyen-Hoang-Nam/let-shorten/form"
	"github.com/Nguyen-Hoang-Nam/let-shorten/models"
	"github.com/Nguyen-Hoang-Nam/let-shorten/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	accountModel = new(models.Account)
	ctx          = context.Background()
)

func Signup(c *gin.Context) {
	var json form.Register
	var err error

	if err = c.ShouldBindJSON(&json); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else if !utils.ValidateRegister(json.Email, json.Password) {
		c.String(http.StatusBadRequest, "Invalid error or password")
	} else {
		var existAccount *models.Account
		existAccount, err = accountModel.GetAccount(json.Email)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		if existAccount != nil {
			c.String(http.StatusBadRequest, "Account exist")
		} else {
			var hashPassword []byte
			hashPassword, err = bcrypt.GenerateFromPassword([]byte(json.Password), 10)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
			} else {
				newAccount := models.Account{
					Email:    json.Email,
					Password: string(hashPassword),
					ID:       uuid.NewString(),
				}

				if err = accountModel.CreateAccount(newAccount); err != nil {
					c.String(http.StatusBadRequest, err.Error())
				} else {
					c.String(http.StatusOK, "Signup successful")
				}
			}
		}
	}
}

func Login(c *gin.Context) {
	var err error
	var json form.Register
	var cookie string

	client := db.GetRedis()

	cookie, err = c.Cookie("gin_cookie")
	if err != nil {
		if err = c.ShouldBindJSON(&json); err != nil {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			var account *models.Account
			account, err = accountModel.GetAccount(json.Email)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
			} else if account == nil {
				c.String(http.StatusNotAcceptable, "Wrong email or password")
			} else if err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(json.Password)); err != nil {
				c.String(http.StatusNotAcceptable, "Wrong email or password")
			} else {
				cookie = uuid.NewString()
				err = client.Set(ctx, cookie, account.ID, 0).Err()
				if err != nil {
					c.String(http.StatusInternalServerError, err.Error())
				} else {
					c.SetCookie("gin_cookie", cookie, 3600, "/", "localhost", false, true)
					c.String(http.StatusOK, "Login successful")
				}
			}
		}
	} else {
		var id string
		id, err = client.Get(ctx, cookie).Result()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			userUrl := fmt.Sprintf("/user/%s", id)
			c.String(http.StatusOK, userUrl)
		}
	}
}

func Signout(c *gin.Context) {
	client := db.GetRedis()

	cookie, err := c.Cookie("gin_cookie")
	if err != nil {
		c.String(http.StatusNotFound, "Not login user")
	} else {
		err = client.Del(ctx, cookie).Err()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusOK, "Sign out successful")
		}
	}
}
