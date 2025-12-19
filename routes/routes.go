package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nerokome/econo/controllers"
	"github.com/nerokome/econo/middleware"
)

// UserRoutes registers all routes for the app
func UserRoutes(router *gin.Engine, app *controllers.Application) {

	// Public routes 
	public := router.Group("/api")
	{
		public.POST("/users/signup", app.SignUp())
		public.POST("/users/login", app.Login())
		public.GET("/users/productview", app.SearchProduct())
		public.GET("/users/search", app.SearchProductByQuery())
	}

	// Protected routes 
	protected := router.Group("/api")
	protected.Use(middleware.Authenticate())
	{
		// Cart
		protected.POST("/cart/add", app.AddToCart())
		protected.GET("/cart/items", app.GetItemFromCart())
		protected.POST("/cart/buy", app.BuyFromCart())
		protected.POST("/cart/instantbuy", app.InstantBuy())

		// Address
		protected.POST("/address", app.AddAddress())
		protected.PUT("/address/:address_id", app.EditAddress())
		protected.DELETE("/address/:address_id", app.DeleteAddress())
	}

	// Admin routes
	admin := router.Group("/admin")
	admin.Use(middleware.Authenticate())
	
	{
		admin.POST("/addproducts", app.ProductViewerAdmin())

	}
}
