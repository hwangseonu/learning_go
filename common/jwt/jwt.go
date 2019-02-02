package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/hwangseonu/goBackend/common/models"
	"net/http"
	"os"
	"strings"
	"time"
)

type CustomClaims struct {
	jwt.StandardClaims
	Identity string `json:"identity"`
}

var secret = os.Getenv("jwt-secret")

func (c CustomClaims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}
	user := new(models.User)
	err := user.FindByUsername(c.Identity)

	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("could not find user")
	}
	return nil
}

func GenerateToken(t, username string) (string, error) {
	var expire int64

	if t == "access" {
		expire = time.Now().Add(time.Hour).Unix()
	} else {
		expire = time.Now().AddDate(0, 1, 0).Unix()
	}

	claims := jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: expire,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "",
		Subject: t,
		NotBefore: time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS512, CustomClaims{claims, username}).SignedString([]byte(secret))
}

func AuthRequire(res http.ResponseWriter, req *http.Request, subject string) *CustomClaims {
	tokenString := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		res.WriteHeader(422)
		res.Write([]byte(`{"message": "jwt error `+err.Error()+`"}`))
		return nil
	}
	claims := token.Claims.(*CustomClaims)

	if claims.Subject != subject {
		res.WriteHeader(422)
		res.Write([]byte(`{"message": "subject required `+subject+`"}`))
		return nil
	}
	return claims
}
