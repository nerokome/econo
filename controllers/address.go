package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nerokome/econo/database"
	"github.com/nerokome/econo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userCollection = database.UserData(database.Client, "users")

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var address models.Address
		if err := c.BindJSON(&address); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		address.ID = primitive.NewObjectID()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := userCollection.UpdateOne(
			ctx,
			bson.M{"user_id": userID},
			bson.M{"$push": bson.M{"addresses": address}},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add address"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "address added"})
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		addressID := c.Param("address_id")

		objID, err := primitive.ObjectIDFromHex(addressID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address id"})
			return
		}

		var updated models.Address
		if err := c.BindJSON(&updated); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := userCollection.UpdateOne(
			ctx,
			bson.M{
				"user_id":       userID,
				"addresses._id": objID,
			},
			bson.M{
				"$set": bson.M{
					"addresses.$.street":  updated.Street,
					"addresses.$.city":    updated.City,
					"addresses.$.pincode": updated.Pincode,
					"addresses.$.house":   updated.House,
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

func EditWorkAddress() gin.HandlerFunc {
	// This naming is cosmetic nonsense — same logic
	return EditHomeAddress()
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		addressID := c.Param("address_id")

		objID, err := primitive.ObjectIDFromHex(addressID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address id"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err = userCollection.UpdateOne(
			ctx,
			bson.M{"user_id": userID},
			bson.M{"$pull": bson.M{"addresses": bson.M{"_id": objID}}},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete address"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "address deleted"})
	}
}
