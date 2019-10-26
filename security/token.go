package security

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var headerTokenKey = "token"

type Token struct {
	ID          uuid.UUID
	AccessLevel int
	Username    string
}

func (t Token) Valid() error {
	if t.AccessLevel < 0 {
		panic("Illegal accessLevel")
	}
	return nil
}

func RetrieveTokenFromRequest(c *gin.Context) (Token, error) {
	tokenString := c.GetHeader(headerTokenKey)
	token, _ := jwt.Parse(tokenString, nil)
	if token == nil || token.Claims == nil {
		return Token{}, errors.New("error while parsing token")
	}
	mapClaim := token.Claims.(jwt.MapClaims)
	uid, _ := uuid.Parse(mapClaim["ID"].(string))
	accessLevel := int(mapClaim["AccessLevel"].(float64))
	securityToken := Token{
		ID:          uid,
		AccessLevel: accessLevel,
		Username:    mapClaim["Username"].(string),
	}
	return securityToken, nil
}
