package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nerokome/econo/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddToCart adds a product to a user's cart
func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Read query params
		productIDStr := c.Query("product_id")
		userIDStr := c.Query("user_id")

		// Validate input
		if productIDStr == "" || userIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "product_id and user_id are required",
			})
			return
		}

		// Convert product ID
		productID, err := primitive.ObjectIDFromHex(productIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid product_id",
			})
			return
		}

		// Convert user ID
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid user_id",
			})
			return
		}

		// Context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// DB call
		err = database.AddProductToCart(
			ctx,
			app.UserCollection,
			app.ProdCollection,
			productID,
			userID,
		)

		if err != nil {
			log.Println("AddToCart error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "product added to cart successfully",
		})
	}
}

// RemoveItem removes a product from cart (stub)
func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"message": "RemoveItem not implemented",
		})
	}
}

// GetItemFromCart returns cart items (stub)
func (app *Application) GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"message": "GetItemFromCart not implemented",
		})
	}
}

// BuyFromCart converts cart to order (stub)
func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"message": "BuyFromCart not implemented",
		})
	}
}

// InstantBuy buys product directly (stub)
func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"message": "InstantBuy not implemented",
		})
	}
}
