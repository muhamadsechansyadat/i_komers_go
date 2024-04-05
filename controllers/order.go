package controllers

import (
	"fmt"
	"i_komers_go/models"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	midtrans "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func AddOrderFromSelectedItemsHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	tx := db.Begin()

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": true, "message": "User data not found in context"})
		return
	}

	userData, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to parse user data"})
		return
	}

	var selectedItems []uint
	if err := c.BindJSON(&selectedItems); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": true, "message": "Invalid request body"})
		return
	}

	var carts []models.Cart
	if err := db.Where("user_id = ? AND (status = '' OR status IS NULL) AND id IN (?)", userData.ID, selectedItems).Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to retrieve user's selected carts"})
		return
	}

	order := models.Order{
		UserID: userData.ID,
		User:   *userData,
		Items:  make([]models.OrderItem, len(carts)),
	}

	for i, cart := range carts {
		var product models.Product
		if err := db.First(&product, cart.ProductID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to retrieve product data"})
			return
		}

		var size models.Size
		if err := db.First(&size, cart.SizeID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to retrieve size data"})
			return
		}

		order.Items[i] = models.OrderItem{
			ProductID: product.ID,
			Product:   product,
			SizeID:    size.ID,
			Size:      size,
			Quantity:  cart.Quantity,
			Price:     calculatePrice(product.Price, cart.Quantity),
		}
		order.Total += order.Items[i].Price
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to create order"})
		return
	}

	if err := tx.Model(&models.Cart{}).Where("user_id = ? AND (status = '' OR statuss IS NULL) AND id IN (?)", userData.ID, selectedItems).Update("status", "N").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to update selected carts status"})
		return
	}

	tx.Commit()

	paymentURL, paymentToken, err := processPayment(&order, userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to process payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"message":     "Order successfully created from selected items",
		"data":        order,
		"payment_url": paymentURL,
		"token":       paymentToken,
	})
}

func processPayment(order *models.Order, userData *models.User) (string, string, error) {
	var snapClient = snap.Client{}
	snapClient.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
	items := []midtrans.ItemDetails{}
	var totalAmount int64
	for _, item := range order.Items {
		items = append(items, midtrans.ItemDetails{
			ID:    strconv.Itoa(int(item.ProductID)),
			Price: int64(item.Product.Price),
			Qty:   int32(item.Quantity),
			Name:  item.Product.Name,
		})
		totalAmount += int64(item.Product.Price) * int64(item.Quantity)
		fmt.Println(item.Product.Price)
		fmt.Println(item.Quantity)
		// fmt.Println()
	}
	fmt.Println(totalAmount)
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order.OrderNumber,
			GrossAmt: totalAmount,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		CustomerDetail: &midtrans.CustomerDetails{
			FName: "Muhamad Sechan",
			LName: "Syadat",
			Email: userData.Email,
			Phone: "087865654343",
		},
		Items: &items,
	}

	chargeResp, err := snapClient.CreateTransaction(req)
	if err != nil {
		return "", "", err
	}

	return chargeResp.RedirectURL, chargeResp.Token, nil
}

func calculatePrice(price float64, quantity uint) float64 {
	return price * float64(quantity)
}
