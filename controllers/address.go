package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nerokome/econo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *Application) AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var address models.Address
		if err := c.ShouldBindJSON(&address); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		address.ID = primitive.NewObjectID()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := app.UserCollection.UpdateOne(
			ctx,
			bson.M{"user_id": userID},
			bson.M{"$push": bson.M{"address_details": address}},
		)

		if err != nil || result.MatchedCount == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add address"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "address added"})
	}
}

func (app *Application) EditAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		addressID := c.Param("address_id")
		objID, err := primitive.ObjectIDFromHex(addressID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address id"})
			return
		}

		var updated models.Address
		if err := c.ShouldBindJSON(&updated); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := app.UserCollection.UpdateOne(
			ctx,
			bson.M{
				"user_id":             userID,
				"address_details._id": objID,
			},
			bson.M{
				"$set": bson.M{
					"address_details.$.street":  updated.Street,
					"address_details.$.city":    updated.City,
					"address_details.$.pincode": updated.Pincode,
					"address_details.$.house":   updated.House,
				},
			},
		)

		if err != nil || result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "address updated"})
	}
}
func (app *Application) DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		addressID := c.Param("address_id")
		objID, err := primitive.ObjectIDFromHex(addressID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := app.UserCollection.UpdateOne(
			ctx,
			bson.M{"user_id": userID},
			bson.M{"$pull": bson.M{
				"address_details": bson.M{"_id": objID},
			}},
		)

		if err != nil || result.MatchedCount == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete address"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "address deleted"})
	}
}
