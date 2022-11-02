package main

import (
	"fmt"
    "net/http"
	"strings"
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

func postLogin(c *gin.Context) {
	// TODO if this is the right thing to do, move to middleware and move url to env var.
	c.Header("Access-Control-Allow-Origin", "http://localhost:4200")
    c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")

	var newUser user

	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		println(err.Error())
		return
	}

	if newUser.Username == "" || newUser.Password == "" || newUser.Token == "" {
		c.JSON(http.StatusBadRequest, "Username, password, and token are required")
		return
	}

	if !strings.Contains(newUser.Username, "@") {
		c.JSON(http.StatusBadRequest, "Usernames must be an email address")
		return
	}

	correct_password, user_exists := users[newUser.Username]

	hour, minutes, _ := time.Now().UTC().Clock()
	correct_token := fmt.Sprintf("%02d%02d", hour, minutes)

	if newUser.Token != correct_token {
		c.JSON(http.StatusUnauthorized, "Token is incorrect")
		return
	}

	// TODO send an array of parameter errors so the frontend can map an error to an input field
	// {"errors":
	// 	[
	// 		{"parameter":"username", "message":"Username not found"},
	// 		{"parameter":"password", "message":"Password is incorrect"},
	// 	]
	// }

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
