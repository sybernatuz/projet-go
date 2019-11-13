package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thoas/go-funk"
	"net/http"
	"time"

	"projet-go/database"
	"projet-go/entities"
	"projet-go/security"
)

func CreateVote(c *gin.Context) {
	userClaims := security.GetUserAuthFromContext(c)
	if userClaims.AccessLevel != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var vote entities.Vote
	err := c.BindJSON(&vote)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	vote.UUID = uuid.New()
	vote.CreatedAt = time.Now()

	if err := database.DBCon.Create(&vote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"uuid":        vote.UUID,
		"title":       vote.Title,
		"description": vote.Description,
	})
}

func GetVote(c *gin.Context) {
	uuidVote, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	vote := entities.Vote{UUID: uuidVote}

	if err := database.DBCon.Where(&vote).First(&vote).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Vote not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"uuid":        vote.UUID,
			"title":       vote.Title,
			"desc": 	   vote.Description,
			"uuid_votes":  vote.UuidVotes,
		})
	}
}

func GetAllVotes(c *gin.Context) {
	var votes []entities.Vote
	database.DBCon.Find(&votes)
	c.JSON(http.StatusOK, gin.H{"votes": votes})
}

func EditVote(c *gin.Context) {
	uuidVote, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	vote := entities.Vote{UUID: uuidVote}

	if err := database.DBCon.Where(&vote).First(&vote).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Vote not found"})
		return
	}

	userClaims := security.GetUserAuthFromContext(c)
	if userClaims.AccessLevel == 1 {
		var voteEdit entities.VoteEdition
		err := c.BindJSON(&voteEdit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}

		vote.UpdatedAt = time.Now()

		if err := database.DBCon.Model(&vote).UpdateColumns(voteEdit).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal error"})
		} else {
			c.JSON(http.StatusOK, vote)
		}
	} else if userClaims.AccessLevel == 0 {
		if funk.Contains(vote.UuidVotes, userClaims.Uuid.String()) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Already voted"})
			return
		}

		uuidVotesUpdated := append(vote.UuidVotes, userClaims.Uuid.String())
		if err := database.DBCon.Model(&vote).Updates(entities.Vote{UuidVotes: uuidVotesUpdated, UpdatedAt: time.Now()}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal error"})
		} else {
			c.JSON(http.StatusOK, vote)
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
	}
}

func DeleteVote(c *gin.Context) {
	userClaims := security.GetUserAuthFromContext(c)
	if userClaims.AccessLevel == 1 {
		uuidVote, err := uuid.Parse(c.Param("uuid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}

		if err := database.DBCon.Delete(&entities.Vote{UUID: uuidVote}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Vote deleted"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
	}
}
