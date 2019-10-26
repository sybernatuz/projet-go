package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"projetgo/database"
	"projetgo/entities"
	"projetgo/security"
)

var tokenHashKey = "secret"

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	user := &entities.User{
		Username: username,
	}
	database.DBCon.Where(user).First(&user)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid login",
		})
		return
	}
	token := security.Token{
		ID:          user.Uuid,
		AccessLevel: user.AccessLevel,
		Username:    user.Username,
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), token)
	tokenString, _ := jwtToken.SignedString([]byte(tokenHashKey))
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Login success",
		"jwt":     tokenString,
	})
}

func Authenticate(c *gin.Context) {
	if !isUriNeedAuthentication(c) {
		return
	}
	token, _ := security.RetrieveTokenFromRequest(c)
	user := entities.User{
		Uuid:        token.ID,
		Username:    token.Username,
		AccessLevel: token.AccessLevel,
	}
	if isUserNotValid(user) {
		errorTokenInvalid(c)
		return
	}
	isTokenNotValid := database.DBCon.Where(&user).First(&user).RecordNotFound()
	if isTokenNotValid {
		errorTokenInvalid(c)
	}
}

func isUriNeedAuthentication(c *gin.Context) bool {
	uri := c.Request.RequestURI
	return uri != "/login" &&
		(uri != "/users/" || c.Request.Method == "post")
}

func isUserNotValid(user entities.User) bool {
	return uuid.Nil == user.Uuid ||
		user.AccessLevel == 0 ||
		user.Username == ""
}

func errorTokenInvalid(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Bad Token",
	})
	c.Abort()
}
