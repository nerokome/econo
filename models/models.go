package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName      string             `json:"first_name" bson:"first_name"`
	LastName       string             `json:"last_name" bson:"last_name"`
	Email          string             `json:"email" bson:"email"`
	Password       string             `json:"password" bson:"password"`
	Tokens         []string           `json:"tokens" bson:"tokens"`
	RefreshTokens  []string           `json:"refresh_tokens" bson:"refresh_tokens"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
	UserID         string             `json:"user_id" bson:"user_id"`
	UserCart       []ProductUser      `json:"user_cart" bson:"user_cart"`
	AddressDetails []Address          `json:"address_details" bson:"address_details"`
	OrderStatus    []Order            `json:"order_status" bson:"order_status"`
}

type Product struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"product_name" bson:"product_name"`
	Price    uint64             `json:"price" bson:"price"`
	Rating   uint8              `json:"rating" bson:"rating"`
	ImageURL string             `json:"image_url" bson:"image_url"`
}

type ProductUser struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"product_name" bson:"product_name"`
	Price    uint64             `json:"price" bson:"price"`
	Rating   uint8              `json:"rating" bson:"rating"`
	ImageURL string             `json:"image_url" bson:"image_url"`
}

type Address struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Street  string             `json:"street" bson:"street"`
	City    string             `json:"city" bson:"city"`
	Pincode string             `json:"pincode" bson:"pincode"`
	House   string             `json:"house" bson:"house"`
}

type Order struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OrderCart   []ProductUser      `json:"order_cart" bson:"order_cart"`
	OrderedAt   time.Time          `json:"ordered_at" bson:"ordered_at"`
	Price       float64            `json:"price" bson:"price"`
	Discount    float64            `json:"discount" bson:"discount"`
	PaymentMode string             `json:"payment_mode" bson:"payment_mode"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	COD     bool `json:"cod" bson:"cod"`
}
