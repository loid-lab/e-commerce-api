package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/loid-lab/e-commerce-api/initializers"
	"github.com/loid-lab/e-commerce-api/models"
	"github.com/loid-lab/e-commerce-api/utils"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var userInput models.User

	recaptchaToken := c.PostForm("recaptchaToken")
	if err := utils.VerifyRecaptcha(recaptchaToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid reCAPTCHA"})
		return
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var userFound models.User
	initializers.DB.Where("email=?", userInput.Email).First(&userFound)

	if userFound.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email is already in use",
		})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := models.User{
		Email:    userInput.Email,
		Password: string(passwordHash),
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create User",
		})
		return
	}

	utils.SendMail(models.EmailData{
		To:       user.Email,
		Subject:  "Welcome to (insert company name here)",
		HTMLBody: fmt.Sprintf("<h2>Welcome %s!</h2><p>Your account has been created</p>", user.Email),
	})
}

func Login(c *gin.Context) {

	recaptchaToken := c.PostForm("recaptchaToken")
	if err := utils.VerifyRecaptcha(recaptchaToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid reCAPTCHA"})
		return
	}

	var userInput models.User

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var userFound models.User
	initializers.DB.Where("email=?", userInput.Email).First(&userFound)

	if userFound.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(userInput.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenStr,
	})
}
