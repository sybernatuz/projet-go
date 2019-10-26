package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"projetgo/database"
	"projetgo/entities"
	"strconv"
)

func DeleteUser(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	database.DBCon.Delete(&entities.User{Uuid: id})
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func GetUser(c *gin.Context) {
	uid, _ := uuid.Parse(c.Query("id"))
	user := entities.User{
		Uuid: uid,
	}
	database.DBCon.Where(&user).First(&user)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GetAllUsers(c *gin.Context) {
	var users []entities.User
	database.DBCon.Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func CreateUser(c *gin.Context) {
	postedPassword := []byte(c.PostForm("password"))
	password, _ := bcrypt.GenerateFromPassword(postedPassword, bcrypt.DefaultCost)
	accessLevel, _ := strconv.Atoi(c.PostForm("AccessLevel"))
	user := entities.User{
		Uuid:        uuid.New(),
		Username:    c.PostForm("username"),
		Password:    string(password),
		FirstName:   c.PostForm("firstName"),
		LastName:    c.PostForm("lastName"),
		Email:       c.PostForm("email"),
		BirthDate:   c.PostForm("birthDate"),
		AccessLevel: accessLevel,
	}
	userDb := database.DBCon.Create(&user)
	c.JSON(http.StatusOK, gin.H{
		"message": "User created",
		"user":    userDb,
	})
}
