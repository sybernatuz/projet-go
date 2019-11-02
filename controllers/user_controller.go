package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"projet-go/database"
	"projet-go/entities"
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
	var user entities.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	user.Uuid = uuid.New()
	password := []byte(user.Password)
	passwordHashed, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	user.Password = string(passwordHashed)

	userDb := database.DBCon.Create(&user)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
		"user":    userDb,
	})
}
