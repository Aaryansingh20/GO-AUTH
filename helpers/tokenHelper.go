package helpers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Aaryansingh20/jwt/database"
	"github.com/Aaryansingh20/jwt/models"
	jwt "github.com/dgrijalva/jwt-go" // golang driver for jwt
	"gorm.io/gorm"
)

type SignedDetails struct {
    Email      string
    First_name string
    Last_name  string
    Uid        string
    User_type  string
    jwt.StandardClaims
}

var userDB *gorm.DB = database.Client

// btw we should have our secret key in .env for production 
var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string) (signedToken string, signedRefreshToken string, err error) {
    claims := &SignedDetails{
        Email:      email,
        First_name: firstName,
        Last_name:  lastName,
        Uid:        uid,
        User_type:  userType,
        StandardClaims: jwt.StandardClaims{
            // setting the expiry time
            ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(120)).Unix(),
        },
    }
    // refreshClaims is used to get a new token if the previous one is expired.

    refreshClaims := &SignedDetails{
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(172)).Unix(),
        },
    }
    token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
    refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
    if err != nil {
        log.Panic(err)
        return
    }
    return token, refreshToken, err
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
    Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
    
    err := userDB.Model(&models.User{}).
        Where("user_id = ?", userId).
        Updates(map[string]interface{}{
            "token":         signedToken,
            "refresh_token": signedRefreshToken,
            "updated_at":    Updated_at,
        }).Error
    
    if err != nil {
        log.Panic(err)
        return
    }
    return
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
    // this function is basically returning the token
    token, err := jwt.ParseWithClaims(
        signedToken,
        &SignedDetails{},
        func(token *jwt.Token) (interface{}, error) {
            return []byte(SECRET_KEY), nil
        },
    )
    if err != nil {
        msg = err.Error()
        return
    }
    // checking if the token is correct or not
    claims, ok := token.Claims.(*SignedDetails)
    if !ok {
        msg = fmt.Sprintf("the token is invalid")
        msg = err.Error()
        return
    }
    // if the token is expired, give error message
    if claims.ExpiresAt < time.Now().Local().Unix() {
        msg = fmt.Sprintf("token has been expired")
        msg = err.Error()
        return
    }
    return claims, msg
}
