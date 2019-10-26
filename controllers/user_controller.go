package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"projetgo/database"
	"projetgo/entities"
	"projetgo/security"
)

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	user := &entities.User{
		Username: username,
	}
	database.DBCon.First(user)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid login",
		})
	}
	token := security.Token{
		ID:          user.Uuid,
		AccessLevel: user.AccessLevel,
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), token)
	tokenString, _ := jwtToken.SignedString([]byte(os.Getenv("token_password")))
	c.Request.Response.Header.Add("token", tokenString)
}

func DeleteUser(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	database.DBCon.Delete(&entities.User{Uuid: id})
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func GetUser(c *gin.Context) {
	user, _ := database.DBCon.Get(c.Query("id"))
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GetAllUsers(c *gin.Context) {
	users := database.DBCon.Find(&entities.User{})
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func CreateUser(c *gin.Context) {
	postedPassword := []byte(c.PostForm("password"))
	password, _ := bcrypt.GenerateFromPassword(postedPassword, bcrypt.DefaultCost)
	user := entities.User{
		Uuid:      uuid.New(),
		Username:  c.PostForm("username"),
		Password:  string(password),
		FirstName: c.PostForm("firstName"),
		LastName:  c.PostForm("lastName"),
		Email:     c.PostForm("email"),
		BirthDate: c.PostForm("birthDate"),
	}
	userDb := database.DBCon.Create(&user)
	c.JSON(http.StatusOK, gin.H{
		"message": "User created",
		"user":    userDb,
	})
}
