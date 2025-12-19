package controllers

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Application holds all shared dependencies for controllers
type Application struct {
	UserCollection *mongo.Collection
	ProdCollection *mongo.Collection
	// Future-proofing:
	// TokenService TokenService
	// Logger       *zap.Logger
	// Config       *config.Config
}

// NewApplication wires dependencies once at startup
func NewApplication(
	userColl *mongo.Collection,
	prodColl *mongo.Collection,
) *Application {
	return &Application{
		UserCollection: userColl,
		ProdCollection: prodColl,
	}
}
