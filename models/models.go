package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json: "_id" bson:"_id"`
	FirstName       string             `bson:"first_name"`
	LastName        string             `bson:"last_name"`
	Email           string             `bson:"email"`
	Password        string             `bson:"password"`
	Tokens          []string           `bson:"tokens"`
	RefreshTokens   []string           `bson:"refresh_tokens"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
	User_ID         string             `bson:"user_id"`
	UserCart        []ProductUser      `bson:"user_cart"`
	Address_Details []Address          `bson:"address_details"`
	Order_Status    []Order            `bson:"order_status"`
}

type Product struct {
	Product_ID   primitive.ObjectID `bson:"product_id"`
	Product_Name string             `bson:"product_name"`
	Price        uint64             `bson:"price"`
	Rating       uint8              `bson:"rating"`
	Image_URL    string             `bson:"image_url"`
}

type ProductUser struct {
	Product_ID   primitive.ObjectID `bson:"product_id"`
	Product_Name string             `bson:"product_name"`
	Price        int                `bson:"price"`
	Rating       uint8              `bson:"rating"`
	Image_URL    string             `bson:"image_url"`
}

type Address struct {
	Address_ID primitive.ObjectID `bson:"address_id"`
	Street     string             `bson:"street"`
	City       string             `bson:"city"`
	Pincode    string             `bson:"pincode"`
	House      string             `bson:"house"`
}

type Order struct {
	Order_ID     primitive.ObjectID `bson:"order_id"`
	Order_cart   []ProductUser       `bson:"order_cart"`
	Ordered_At   time.Time          `bson:"ordered_at"`
	Price        float64            `bson:"price"`
	Discount     float64            `bson:"discount"`
	Payment_Mode string             `bson:"payment_mode"`
}

type Payment struct {
	Digital bool `bson:"digital"`
	COD     bool `bson:"cod"`
}
