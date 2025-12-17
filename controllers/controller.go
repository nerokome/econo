package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nerokome/econo/database"
	"github.com/nerokome/econo/models"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

// HashPassword hashes a plain-text password
func HashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashed)
}

// VerifyPassword compares stored hash with incoming password
func VerifyPassword(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	)
	return err == nil
}

// SIGN UP
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// basic validation (email & password should not be empty)
		if user.Email == "" || user.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "email and password are required",
			})
			return
		}

		collection := database.UserData(database.Client, "Users")

		// check email uniqueness
		count, err := collection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}

		// hash password
		user.Password = HashPassword(user.Password)

		// metadata
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		user.Tokens = []string{}
		user.RefreshTokens = []string{}
		user.UserCart = []models.ProductUser{}
		user.AddressDetails = []models.Address{}
		user.OrderStatus = []models.Order{}

		_, err = collection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user creation failed"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "user created successfully",
		})
	}
}

// LOGIN
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var loginData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var foundUser models.User

		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		collection := database.UserData(database.Client, "Users")

		err := collection.
			FindOne(ctx, bson.M{"email": loginData.Email}).
			Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		if !VerifyPassword(foundUser.Password, loginData.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "login successful",
			"user":    foundUser,
		})
	}
}

// STUBS (compile-safe)

// ProductViewerAdmin returns all products for admin view
func (app *Application) ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := app.ProdCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
			return
		}
		defer cursor.Close(ctx)

		var products []models.Product
		if err := cursor.All(ctx, &products); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode products"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"products": products})
	}
}

// SearchProduct returns all products (public)
func (app *Application) SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := app.ProdCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
			return
		}
		defer cursor.Close(ctx)

		var products []models.ProductUser // public view can hide internal fields
		for cursor.Next(ctx) {
			var p models.Product
			if err := cursor.Decode(&p); err != nil {
				continue
			}
			products = append(products, models.ProductUser{
				ID:       p.ID,
				Name:     p.Name,
				Price:    p.Price,
				Rating:   p.Rating,
				ImageURL: p.ImageURL,
			})
		}

		c.JSON(http.StatusOK, gin.H{"products": products})
	}
}

// SearchProductByQuery searches products by name query
func (app *Application) SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'q' is required"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.M{"product_name": bson.M{"$regex": query, "$options": "i"}}

		cursor, err := app.ProdCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search products"})
			return
		}
		defer cursor.Close(ctx)

		var products []models.ProductUser
		for cursor.Next(ctx) {
			var p models.Product
			if err := cursor.Decode(&p); err != nil {
				continue
			}
			products = append(products, models.ProductUser{
				ID:       p.ID,
				Name:     p.Name,
				Price:    p.Price,
				Rating:   p.Rating,
				ImageURL: p.ImageURL,
			})
		}

		c.JSON(http.StatusOK, gin.H{"products": products})
	}
}

