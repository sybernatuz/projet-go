package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"projet-go/controllers"
	"projet-go/database"
	"projet-go/entities"
)

func main() {
	database.DBCon, database.Error = gorm.Open("postgres", "host=db port=5432 user=user dbname=api-vote password=password sslmode=disable")

	if database.Error != nil {
		panic(database.Error)
	}

	database.DBCon.AutoMigrate(&entities.User{})
	database.DBCon.AutoMigrate(&entities.Vote{})
	database.DBCon.AutoMigrate(&entities.Ip{})
	controllers.Router = gin.Default()

	controllers.Router.Use(controllers.Authenticate)

	controllers.Router.NoRoute(controllers.ErrorNotFound)

	controllers.Router.GET("/", controllers.Home)

	controllers.Router.POST("/login", controllers.Login)

	controllers.Router.GET("/users/:uuid", controllers.GetUser)
	controllers.Router.GET("/users/", controllers.GetAllUsers)
	controllers.Router.POST("/users/", controllers.CreateUser)
	controllers.Router.DELETE("/users/:uuid", controllers.DeleteUser)
	controllers.Router.PUT("/users/:uuid", controllers.UpdateUser)

	controllers.Router.POST("/votes/", controllers.CreateVote)
	controllers.Router.GET("/votes/:uuid", controllers.GetVote)
	controllers.Router.GET("/votes/", controllers.GetAllVotes)
	controllers.Router.PUT("/votes/:uuid", controllers.EditVote)
	controllers.Router.DELETE("/votes/:uuid", controllers.DeleteVote)

	_ = controllers.Router.Run()
}
