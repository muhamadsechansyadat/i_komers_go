package controllers

import (
	"i_komers_go/helpers"
	"i_komers_go/middleware"
	"i_komers_go/models"
	"net/http"

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

	db.Create(&user)
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
