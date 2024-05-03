package controllers

import (
	"fmt"
	"i_komers_go/helpers"
	"i_komers_go/middleware"
	"i_komers_go/models"
	"net/http"
	"net/smtp"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	var user models.User

	var registerRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}

	if err := c.ShouldBind(&registerRequest); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			errors := helpers.ParseError(validationErr)
			helpers.HelperErrorWithDataJSON(c, errors)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	db.Where("username = ?", registerRequest.Username).First(&user)
	if user.ID != 0 {
		helpers.ErrorJSON(c, "Username already taken.")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.ErrorJSON(c, "Failed to hash password.")
		return
	}

	user.Username = registerRequest.Username
	user.Password = string(hashedPassword)
	user.Email = registerRequest.Email

	// db.Create(&user)

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to convert SMTP_PORT to integer"})
		return
	}
	emailAddress := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	from := emailAddress
	to := []string{"muhamadsechansyadat@gmail.com"}
	subject := "Test Email"
	body := "Hello, this is a test email from Go!"

	message := fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("To: %s\r\n", to[0]) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" +
		body

	auth := smtp.PlainAuth("", emailAddress, password, smtpHost)
	err = smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, from, to, []byte(message))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Created successfully"})
}

func LoginHandler(c *gin.Context) {
	var user models.User

	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		helpers.ErrorJSON(c, err.Error())
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	db.Where("username = ?", loginRequest.Username).First(&user)
	if user.ID == 0 {
		helpers.ErrorJSON(c, "invalid credential 1.")
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		helpers.ErrorJSON(c, "Password does not match.")
		return
	}

	token, err := middleware.GenerateToken(user)
	if err != nil {
		helpers.ErrorJSON(c, "Error generating token.")
		return
	}

	helpers.SuccessJSON(c, "login successfully", token)
}
