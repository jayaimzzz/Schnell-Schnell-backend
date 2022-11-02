package main

import (
	"fmt"
    "net/http"
	"time"
	
	"golang.org/x/crypto/bcrypt"
    "github.com/gin-gonic/gin"
)

type user struct {
	Username	string 	`json:"username"`
	Password	string	`json:"password"`
	Token 		string	`json:"token"`
}

var users = map[string]string{
	"c137@onecause.com": "$2a$10$rjcgEIq4.DFlTnXbSMwgOeBO84VagAt6GkxgdYTD/lBLGFawDwtZ6",
}

func main() {
    router := gin.Default()
    router.POST("/login", postLogin)

    router.Run("localhost:8080")
}

// func HashPassword(password string) (string, error) {
//     bytes, err := bcrypt.GenerateFromPassword([]byte(password), 1)
//     return string(bytes), err
// }

func postLogin(c *gin.Context) {
	var newUser user

	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if newUser.Username == "" || newUser.Password == "" || newUser.Token == "" {
		c.JSON(http.StatusBadRequest, "Username, password, and token are required")
		return
	}

	correct_password, user_exists := users[newUser.Username]

	hour, minutes, _ := time.Now().UTC().Clock()
	correct_token := fmt.Sprintf("%02d%02d", hour, minutes)

	if newUser.Token != correct_token {
		c.JSON(http.StatusUnauthorized, "Token is incorrect")
		return
	}

	if !user_exists {
		c.JSON(http.StatusUnauthorized, "Username not found")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(correct_password), []byte(newUser.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Password is incorrect")
		return
	}

	c.JSON(http.StatusOK, "authorized")
}
