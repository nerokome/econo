package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nerokome/econo/models"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashed)
}

func VerifyPassword(hashedPassword, plainPassword string) bool {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	) == nil
}

func (app *Application) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if user.Email == "" || user.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "email and password are required",
			})
			return
		}

		count, err := app.UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}

		user.Password = HashPassword(user.Password)
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		user.Tokens = []string{}
		user.RefreshTokens = []string{}
		user.UserCart = []models.ProductUser{}
		user.AddressDetails = []models.Address{}
		user.OrderStatus = []models.Order{}

		if _, err := app.UserCollection.InsertOne(ctx, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user creation failed"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "user created successfully",
		})
	}
}

func (app *Application) Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var creds struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		err := app.UserCollection.FindOne(ctx, bson.M{"email": creds.Email}).Decode(&user)
		if err != nil || !VerifyPassword(user.Password, creds.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "login successful",
			"user":    user,
		})
	}
}
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
func (app *Application) SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {

		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'q' is required"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.M{
			"name": bson.M{"$regex": query, "$options": "i"},
		}

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
