package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"projetgo/controllers"
	"projetgo/database"
	"projetgo/entities"
)

func main() {
	database.DBCon, database.Error = gorm.Open("postgres", "host=db port=5432 user=user dbname=api-vote password=password sslmode=disable")

	if database.Error != nil {
		panic(database.Error)
	}
	database.DBCon.AutoMigrate(&entities.User{})
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/users/login", controllers.Login)
	r.GET("/users/:uuid", controllers.GetUser)
	r.GET("/users/", controllers.GetAllUsers)
	r.POST("/users/", controllers.CreateUser)
	r.DELETE("/users/:uuid", controllers.DeleteUser)
	_ = r.Run()
}
