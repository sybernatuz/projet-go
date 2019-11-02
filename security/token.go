package security

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
)

var tokenHashKey = []byte("secret")

type UserClaims struct {
	Uuid        uuid.UUID
	AccessLevel int
	jwt.StandardClaims
}

func (t UserClaims) Valid() error {
	if t.AccessLevel < 0 {
		panic("Illegal accessLevel")
	}
	return nil
}

func JwtDecode(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return tokenHashKey, nil
	})
}

func JwtCreate(uuid uuid.UUID, accessLevel int, expiredAt int64) (string, error) {
	claims := UserClaims{
		uuid,
		accessLevel,
		jwt.StandardClaims{
			ExpiresAt: expiredAt,
			Issuer:    "projet-go",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(tokenHashKey)
	if err != nil {
		return "", errors.New("error with token")
	}
	return ss, nil
}

func RetrieveTokenFromRequest(c *gin.Context) (string, error) {
	reqToken := c.GetHeader("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("invalid token")
	}
	reqToken = splitToken[1]

	return reqToken, nil
}

func UserInfosFromToken(tokenString string) (*UserClaims, error) {
	token, err := JwtDecode(tokenString)
	if err != nil {
		return nil, errors.New("bad token")
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		if claims == nil {
			return nil, errors.New("bad token")
		}
		return claims, nil
	} else {
		return nil, errors.New("bad token")
	}
}

func GetUserAuthFromContext(c *gin.Context) *UserClaims {
	userAuth, _ := c.Get("user")
	userClaims := userAuth.(*UserClaims)
	return userClaims
}
