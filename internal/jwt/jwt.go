package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/janiokq/Useless-blog/cinit"

	"time"
)

var (
	secret = []byte("Your key secret")
)

type Msg struct {
	UserID   int64  `json:"userid"`
	UserName string `json:"username"`
}

type MyClaims struct {
	Msg
	jwt.StandardClaims
}

func Encode(r Msg) (string, error) {
	claims := MyClaims{
		r,
		jwt.StandardClaims{
			// ExpiresAt: time.Now().Unix(),// 过期时间
			ExpiresAt: time.Now().Add(time.Hour * cinit.TokenExpirationtime).Unix(), //  过 期时间
			Issuer:    "live.incode",                                                //  该JWT的签发者,可选
			// Subject:"",// 该JWT所面向的用户
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func Decode(tokenString string) (Msg, error) {
	//  Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return Msg{}, err
	}
	//  Validate the token and return the custom claims
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims.Msg, nil
	}
	return Msg{}, err
}
