package database

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct      = errors.New("cannot find product")
	ErrorCantDecodeProducts = errors.New("cannot decode products")
	ErrProductNotFound      = errors.New("product not found")
	ErrCartEmpty            = errors.New("cart is empty")
	ErrInsufficientQuantity = errors.New("insufficient quantity")
	ErrUnableToUpdateCart   = errors.New("unable to update cart")
	ErrUnableToCreateOrder  = errors.New("unable to create order")
	ErrOrderCreationFailed  = errors.New("order creation failed")
	ErrInstantBuyFailed     = errors.New("instant buy failed")
	ErrUserIdIsnotValid     = errors.New("user id is not valid")
)

/*
AddProductToCart pushes a product ID into user's cart
*/
func AddProductToCart(
	ctx context.Context,
	userCollection *mongo.Collection,
	productCollection *mongo.Collection,
	userID primitive.ObjectID,
	productID primitive.ObjectID,
) error {

	// Check product exists
	count, err := productCollection.CountDocuments(
		ctx,
		bson.M{"_id": productID},
	)
	if err != nil || count == 0 {
		return ErrProductNotFound
	}

	// Add product to cart (no duplicates)
	result, err := userCollection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$addToSet": bson.M{
				"cart": productID,
			},
		},
	)

	if err != nil || result.MatchedCount == 0 {
		return ErrUnableToUpdateCart
	}

	return nil
}

/*
RemoveProductFromCart removes a product from user's cart
*/
func RemoveProductFromCart(
	ctx context.Context,
	userCollection *mongo.Collection,
	userID primitive.ObjectID,
	productID primitive.ObjectID,
) error {

	result, err := userCollection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$pull": bson.M{
				"cart": productID,
			},
		},
	)

	if err != nil || result.MatchedCount == 0 {
		return ErrUnableToUpdateCart
	}

	return nil
}

/*
GetUserCart returns product IDs in cart
*/
func GetUserCart(
	ctx context.Context,
	userCollection *mongo.Collection,
	userID primitive.ObjectID,
) ([]primitive.ObjectID, error) {

	var user struct {
		Cart []primitive.ObjectID `bson:"cart"`
	}

	err := userCollection.FindOne(
		ctx,
		bson.M{"_id": userID},
	).Decode(&user)

	if err != nil {
		return nil, ErrUserIdIsnotValid
	}

	if len(user.Cart) == 0 {
		return nil, ErrCartEmpty
	}

	return user.Cart, nil
}
