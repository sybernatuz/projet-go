package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"net/http"
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
	controllers.Router = gin.Default()

	controllers.Router.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "404",
		})
	})
	controllers.Router.GET("/", controllers.Home)
	controllers.Router.POST("/users/login", controllers.Login)
	controllers.Router.GET("/users/:uuid", controllers.GetUser)
	controllers.Router.GET("/users/", controllers.GetAllUsers)
	controllers.Router.POST("/users/", controllers.CreateUser)
	controllers.Router.DELETE("/users/:uuid", controllers.DeleteUser)
	_ = controllers.Router.Run()
}
