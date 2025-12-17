package controllers

import "go.mongodb.org/mongo-driver/mongo"

// Application is shared across all controllers
type Application struct {
	ProdCollection *mongo.Collection
	UserCollection *mongo.Collection
}

func NewApplication(prodColl, userColl *mongo.Collection) *Application {
	return &Application{
		ProdCollection: prodColl,
		UserCollection: userColl,
	}
}
