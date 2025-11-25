package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	database "github.com/Aaryansingh20/jwt/database"
	helper "github.com/Aaryansingh20/jwt/helpers"
	models "github.com/Aaryansingh20/jwt/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var userDB *gorm.DB = database.Client
var validate = validator.New()

func HashPassword(password string) string {
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        log.Panic(err)
    }
    return string(hashed)
}

func VerifyPassword(userPassword, providedPassword string) (bool, string) {
    err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
    check := true
    msg := ""
    if err != nil {
        check = false
        msg = fmt.Sprintf("email or password is incorrect.")
    }
    return check, msg
}

func SignUp() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()
        
        var user models.User

        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        validationErr := validate.Struct(user)
        // this is used to validate, but what? see the User struct, and see those validate struct fields
        if validationErr != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
            return
        }

        // Check if email already exists (PostgreSQL version)
        var count int64
        userDB.WithContext(ctx).Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
        if count > 0 {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
            return
        }

        password := HashPassword(*user.Password)
        user.Password = &password

        // Check if phone already exists (PostgreSQL version)
        userDB.WithContext(ctx).Model(&models.User{}).Where("phone = ?", user.Phone).Count(&count)
        if count > 0 {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "this phone number already exists"})
            return
        }

        // by "c.BindJSON(&user)" user already have the information from the website user
        user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        user.User_id = fmt.Sprintf("%d", time.Now().UnixNano()) // Generate unique ID
        
        token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)

        // giving value that we generated to user
        user.Token = &token
        user.Refresh_token = &refreshToken

        // now let's insert it to the database (PostgreSQL version)
        result := userDB.WithContext(ctx).Create(&user)
        if result.Error != nil {
            msg := fmt.Sprintf("User item was not created")
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
            return
        }

        c.JSON(http.StatusOK, user)
    }
}

func Login() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()
        
        var user models.User
        var foundUser models.User

        // giving the user data to user variable
        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // finding the user through email (PostgreSQL version)
        err := userDB.WithContext(ctx).Where("email = ?", user.Email).First(&foundUser).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
            return
        }

        // we need pointer to access the original user and foundUser,
        // if we only pass user and foundUser, it will create a new instance of user and foundUser
        isPasswordValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
        if isPasswordValid != true {
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
            return
        }

        if foundUser.Email == nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
            return
        }
        
        token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
        helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
        
        // Get updated user with new tokens (PostgreSQL version)
        err = userDB.WithContext(ctx).Where("user_id = ?", foundUser.User_id).First(&foundUser).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(http.StatusOK, foundUser)
    }
}

// GetUsers can only be accessed by the admin.
func GetUsers() gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := helper.CheckUserType(c, "ADMIN"); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // setting how many records you want per page.
        // we are taking the recordPerPage from c and converting it to int
        recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

        // if error or recordPerPage is less than 1, by default we will have 9 records per page
        if recordPerPage < 1 || err != nil {
            recordPerPage = 9
        }

        // this is just like page number
        page, err1 := strconv.Atoi(c.Query("page"))
        // we want to start with the page number 1 by default.
        if err1 != nil || page < 1 {
            page = 1
        }

        startIndex := (page - 1) * recordPerPage

        var users []models.User
        var totalCount int64

        // Get total count (PostgreSQL version)
        userDB.Model(&models.User{}).Count(&totalCount)

        // Get paginated users (PostgreSQL version)
        result := userDB.Limit(recordPerPage).Offset(startIndex).Find(&users)
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing user items"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "total_count": totalCount,
            "user_items":  users,
        })
    }
}

func GetUserById() gin.HandlerFunc {
    return func(c *gin.Context) {
        userId := c.Param("user_id") // we are taking the user_id given by the user in json
        // with the help of gin.context we can access the json data send by postman or curl or user

        if err := helper.MatchUserTypeToUserId(c, userId); err != nil {
            //checking if the user in admin or not.
            // we will create that func in helper package.

            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        var user models.User

        // Find user by user_id (PostgreSQL version)
        err := userDB.WithContext(ctx).Where("user_id = ?", userId).First(&user).Error
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // if everything goes ok, pass the data of the user (UserModel.go)
        c.JSON(http.StatusOK, user)
    }
}
