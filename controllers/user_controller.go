package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"projet-go/database"
	"projet-go/entities"
	"projet-go/security"
	"time"
	"github.com/thedevsaddam/govalidator"
)

type UserRequestParams struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"pass,omitempty"`
	DateOfBirth string `json:"birth_date,omitempty"`
}

func DeleteUser(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("uuid"))
	database.DBCon.Delete(&entities.User{UUID: id})
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func GetUser(c *gin.Context) {
	uid, _ := uuid.Parse(c.Param("uuid"))
	user := entities.User{
		UUID: uid,
	}
	database.DBCon.Where(&user).First(&user)
	c.JSON(http.StatusOK, user)
}

func GetAllUsers(c *gin.Context) {
	var users []entities.User
	database.DBCon.Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func UpdateUser(c *gin.Context) {
	userClaims := security.GetUserAuthFromContext(c)
	if userClaims.AccessLevel != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	var requestParams UserRequestParams
	err := c.ShouldBindJSON(&requestParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	uid, _ := uuid.Parse(c.Param("uuid"))
	user := entities.User{
		UUID: uid,
	}
	database.DBCon.Where(user).First(&user)
	if requestParams.LastName != "" {
		user.LastName = requestParams.LastName
	}
	if requestParams.FirstName != "" {
		user.FirstName = requestParams.FirstName
	}
	if requestParams.Password != "" {
		user.Password = requestParams.Password
	}
	if requestParams.Email != "" {
		user.Email = requestParams.Email
	}
	if requestParams.DateOfBirth != "" {
		date, err := time.Parse("02-01-2006", requestParams.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad date"})
			return
		}
		user.DateOfBirth = date
	}
	userDb := database.DBCon.Save(user)
	c.JSON(http.StatusOK, userDb.Value)
}

func CreateUser(c *gin.Context) {
	var requestParams UserRequestParams
	rules := govalidator.MapData{
		"first_name": 	   []string{"required", "min:2"},
		"last_name": 	   []string{"required", "min:2"},
		"email":    	   []string{"required", "email"},
		"birth_date":      []string{"url"},
	}

	opts := govalidator.Options{
		Request: c.Request,
		Data:    &UserRequestParams,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	err := c.BindJSON(&requestParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	date, err := time.Parse("02-01-2006", requestParams.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad date"})
		return
	}
	user := entities.User{
		Email:       requestParams.Email,
		Password:    requestParams.Password,
		FirstName:   requestParams.FirstName,
		LastName:    requestParams.LastName,
		DateOfBirth: date,
	}
	user.UUID = uuid.New()
	user.CreatedAt = time.Now()
	password := []byte(user.Password)
	passwordHashed, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	user.Password = string(passwordHashed)

	userDb := database.DBCon.Create(&user)
	if userDb.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while creating user"})
		return
	}

	c.JSON(http.StatusCreated, userDb.Value)
}
