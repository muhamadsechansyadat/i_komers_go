package controllers

import (
	"i_komers_go/helpers"
	"i_komers_go/middleware"
	"i_komers_go/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	var user models.User

	var registerRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		helpers.ErrorJSON(c, err.Error())
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

	db.Create(&user)
	helpers.SuccessStringJSON(c, "User created successfully")
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
		helpers.ErrorJSON(c, "invalid credential 2.")
		return
	}

	token, err := middleware.GenerateToken(user)
	if err != nil {
		helpers.ErrorJSON(c, "Error generating token.")
		return
	}

	helpers.SuccessJSON(c, "login successfully", token)
}
