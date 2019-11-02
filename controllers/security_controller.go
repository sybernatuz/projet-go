package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"projet-go/database"
	"projet-go/entities"
	"projet-go/security"
	"regexp"
	"time"
)

type UserAuth struct {
	Username string
	Password string
}

func Login(c *gin.Context) {
	var userAuth UserAuth
	err := c.BindJSON(&userAuth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	user := &entities.User{
		Username: userAuth.Username,
	}

	database.DBCon.Where(user).First(&user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userAuth.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid login",
		})
		return
	}

	expiredAt := time.Now().Add(time.Hour * 1).Unix()
	tokenString, err := security.JwtCreate(user.Uuid, user.AccessLevel, expiredAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Login success",
		"jwt":       tokenString,
		"expiredAt": expiredAt,
	})
}

func Authenticate(c *gin.Context) {
	if !isUriNeedAuthentication(c) {
		return
	}

	tokenString, err := security.RetrieveTokenFromRequest(c)
	if err != nil {
		errorTokenInvalid(c)
		return
	}

	userAuth, err := security.UserInfosFromToken(tokenString)
	if err != nil {
		errorTokenInvalid(c)
	}

	c.Set("user", userAuth)
}

func isUriNeedAuthentication(c *gin.Context) bool {
	uri := c.Request.RequestURI

	matchVotes, _ := regexp.MatchString("/votes/*", uri)
	matchUsers, _ := regexp.MatchString("/users/*", uri)

	return uri != "/login/" && ((matchUsers && c.Request.Method != "GET" && c.Request.Method != "POST") ||
		(matchVotes && c.Request.Method != "GET"))
}

func errorTokenInvalid(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Bad Token",
	})
	c.Abort()
}
