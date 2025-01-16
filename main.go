package main

import (
	"fmt"
	"go_project/internal/handlers"
	"go_project/internal/models"
	"go_project/middlewares"
	"go_project/pkg/db"
	"go_project/pkg/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Import DB and TestConnection
	db.InitDB()
	db.DB.AutoMigrate(&models.User{}, &models.UserInfo{}, &models.DeliveryAddress{}, &models.PickupPoint{}, &models.Verification{}, &models.Category{}, &models.Item{}, &models.ItemImage{}, &models.RecentlyViewedItem{}, &models.Claims{}, &models.Order{}, &models.OrderItem{},  
	 &models.CartItem{}, &models.LikedItem{})
	r, err := router.InitRouter()

	if err != nil {
		fmt.Errorf(" Error init router %w ", err.Error())
		return
	}

	r.LoadHTMLGlob("static/**/*.html")
	r.Static("/static", "./static/web")
	r.Static("/assets", "./static/assets")
	r.Static("/admin/static", "./static/admin")
	r.Static("/seller/static", "./static/seller")

	//for images
	// Обработка статических файлов из директории uploads
	r.Static("/uploads", "./uploads")

	

	userRepo := handlers.NewUserRepository(db.DB)

	r.Handle("GET", "/user", middlewares.IsAuthorized(), middlewares.IsAdmin(), userRepo.GetAllUsers)
	r.Handle("POST", "/signup", handlers.Signup)
	r.Handle("GET", "/signup", handlers.Signup)
	r.Handle("GET", "/user/:id", middlewares.IsAuthorized(), middlewares.IsAdmin(), userRepo.GetUserByID)
	r.Handle("DELETE", "/user/:id", middlewares.IsAuthorized(), middlewares.IsAdmin(), userRepo.DeleteUser)
	r.Handle("PUT", "/user/:id", middlewares.IsAuthorized(), middlewares.IsAdmin(),  userRepo.UpdateUser)
	r.Handle("POST", "/login", handlers.Login)
	r.Handle("GET", "/login", handlers.Login)

	r.Handle("GET", "/logout", handlers.Logout)
	r.Handle("GET", "/verification", handlers.Verification)
	r.Handle("POST", "/verification", handlers.Verification)
	r.Handle("GET", "/reset-password", handlers.ResetPassword)
	r.Handle("POST", "/reset-password", handlers.ResetPassword)

	r.Handle("GET", "/verify-reset-code", handlers.VerifyResetCode)
	r.Handle("POST", "/verify-reset-code", handlers.VerifyResetCode)

	r.Handle("GET", "/update-password", handlers.UpdatePassword)
	r.Handle("POST", "/update-password", handlers.UpdatePassword)

	
	itemRepo := handlers.NewItemRepository(db.DB)
	cartRepo := handlers.NewCartRepository(db.DB)
	
	

	r.Handle("GET", "/products", middlewares.IsAuthorized(), itemRepo.GetAllItems)
	// r.Handle("GET", "/products/:id", middlewares.IsAuthorized(), itemRepo.GetItemByID)
	r.Handle("POST", "/products", middlewares.IsAuthorized(), middlewares.IsAdmin(), itemRepo.CreateItem)
	r.Handle("PUT", "/products/:id", middlewares.IsAuthorized(), middlewares.IsAdmin(), itemRepo.UpdateItem)
	r.Handle("DELETE", "/products/:id", middlewares.IsAuthorized(), middlewares.IsAdmin(), itemRepo.DeleteItem)

	categoryRepo := handlers.NewCategoryRepository(db.DB)
	r.Handle("GET", "/categories", middlewares.IsAuthorized(), categoryRepo.GetAllCategories)
	r.Handle("POST", "/categories", middlewares.IsAuthorized(), middlewares.IsAdmin(), categoryRepo.CreateCategory)
	r.Handle("GET", "/categories/:id", middlewares.IsAuthorized(), categoryRepo.GetCategoryByID)
	r.Handle("PUT", "/categories/:id", middlewares.IsAuthorized(), middlewares.IsAdmin(), categoryRepo.UpdateCategory)
	r.Handle("DELETE", "/categories/:id", middlewares.IsAuthorized(), middlewares.IsAdmin(), categoryRepo.DeleteCategory)

	orderRepo := handlers.NewOrderRepository(db.DB)
	// r.Handle("POST", "/orders", middlewares.IsAuthorized(), orderRepo.CreateOrder)
	r.Handle("PUT", "/orders/:id", middlewares.IsAuthorized(), middlewares.IsAdmin(), orderRepo.UpdateOrder)
	r.Handle("GET", "/orders/:id", middlewares.IsAuthorized(), middlewares.IsAdmin(), orderRepo.GetOrderByID)
	r.Handle("GET", "/orders", middlewares.IsAuthorized(), middlewares.IsAdmin(), orderRepo.GetAllOrders)
	r.Handle("DELETE", "/orders/:id", middlewares.IsAuthorized(), middlewares.IsAdmin(), orderRepo.DeleteOrder)

	//for seller
	r.Handle("GET", "/seller/dashboard", middlewares.IsAuthorized(), middlewares.SellerOnly(), itemRepo.SellerDashboard)

	r.Handle("GET", "/seller/products", middlewares.IsAuthorized(), middlewares.SellerOnly(), itemRepo.SellerGetAllItems)
	r.Handle("POST", "/seller/products", middlewares.IsAuthorized(), middlewares.SellerOnly(), itemRepo.SellerCreateItem)
	r.Handle("GET", "/seller/create/products", middlewares.IsAuthorized(), middlewares.SellerOnly(), itemRepo.SellerCreateItem)
	r.Handle("GET", "/seller/edit/product/:id", middlewares.IsAuthorized(), middlewares.SellerOnly(), itemRepo.SellerEditItem)
	r.Handle("POST", "/seller/edit/product/:id", middlewares.IsAuthorized(), middlewares.SellerOnly(), itemRepo.SellerUpdateItem)
	r.Handle("DELETE", "/seller/products/delete/:id", middlewares.IsAuthorized(), middlewares.SellerOnly(), itemRepo.SellerDeleteItem)



	//for images
	r.Handle("POST", "/upload", handlers.UploadImage)

	recentlyViewedRepo := handlers.NewRecentlyViewedRepository(db.DB)
	
	r.Handle("POST", "/cart/add", middlewares.IsAuthorized(), handlers.AddToCartHandler(cartRepo))

	r.Handle("DELETE", "/cart/remove", middlewares.IsAuthorized(), handlers.RemoveFromCartHandler(cartRepo))
	r.Handle("POST", "/update-cart-quantity", middlewares.IsAuthorized(), handlers.UpdateCartQuantityHandler(cartRepo))

	likedRepo := handlers.NewLikedRepository(db.DB)
	r.Handle("POST", "/liked/add", middlewares.IsAuthorized(), handlers.AddToLikedHandler(likedRepo))
	r.Handle("GET", "/liked", middlewares.IsAuthorized(), handlers.ViewLikedHandler(likedRepo, cartRepo))
	r.Handle("GET", "/cart", middlewares.IsAuthorized(), handlers.ViewCartHandler(itemRepo, cartRepo, likedRepo, recentlyViewedRepo))


	// for item 


	r.Handle("GET", "/item/:id", middlewares.IsAuthorized(), func(c *gin.Context) {
		itemRepo.GetItemByID(c, recentlyViewedRepo)
	})

	r.Handle("GET", "/map", middlewares.IsAuthorized(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "map.html", nil)
	})
	r.Handle("GET", "/api/pickup-points", middlewares.IsAuthorized(), handlers.GetPickupPoints)

	deliveryAddressRepo := handlers.NewDeliveryAddressRepository(db.DB)
	r.Handle("POST", "/api/save-delivery-address", middlewares.IsAuthorized(), handlers.SaveDeliveryAddressHandler(deliveryAddressRepo))
	r.Handle("GET", "/api/get-delivery-address", middlewares.IsAuthorized(), handlers.GetDeliveryAddressHandler(deliveryAddressRepo))
	
	itemCtrl := handlers.NewItemController(itemRepo)
	r.Handle("GET", "/search", middlewares.IsAuthorized(), itemCtrl.HandleSearchItems)
	r.Handle("GET", "/cart/item-count", handlers.GetCartItemCountHandler(cartRepo))
	r.Handle("GET", "/home",  middlewares.IsAuthorized(), itemRepo.HomePage(cartRepo, likedRepo))

	r.Handle("GET", "/checkout", middlewares.IsAuthorized(), handlers.ViewCheckoutHandler(cartRepo))
	r.Handle("POST", "/create-payment-intent", handlers.CreatePaymentIntentHandler)
	r.Handle("GET", "/create-order", middlewares.IsAuthorized(), handlers.CreateOrder(orderRepo, cartRepo))
	r.Handle("GET", "/my-orders", middlewares.IsAuthorized(), handlers.MyOrdersHandler(orderRepo))
	r.Handle("GET", "/order/:id", middlewares.IsAuthorized(), handlers.GetOrderDetailsHandler(orderRepo))


	r.Handle("GET", "/profile", middlewares.IsAuthorized(), handlers.ProfileHandler(userRepo))
	r.Handle("POST", "/profile-update", middlewares.IsAuthorized(), handlers.UpdateProfileHandler(userRepo))

	adminRepo := handlers.NewAdminRepository(db.DB)
	r.Handle("GET", "/admin", middlewares.IsAuthorized(), middlewares.IsAdmin(), handlers.AdminDashboardHandler(adminRepo))
	r.Handle("GET", "/admin-users", middlewares.IsAuthorized(), middlewares.IsAdmin(), handlers.AdminUsersHandler(adminRepo))
	r.Handle("GET", "/admin/orders", middlewares.IsAuthorized(), middlewares.IsAdmin(), handlers.AdminOrdersHandler(adminRepo))
	r.Handle("GET", "/admin/order-details/:id", middlewares.IsAuthorized(), middlewares.IsAdmin(), handlers.AdminOrderDetailsHandler(adminRepo))
	r.Handle("GET", "/admin/order-update/:id", middlewares.IsAuthorized(), middlewares.IsAdmin(), handlers.AdminOrderUpdateHandler(adminRepo))
	r.Handle("POST", "/admin/order-update/:id", middlewares.IsAuthorized(), middlewares.IsAdmin(), handlers.AdminOrderUpdateHandler(adminRepo))





	reviewRepo := handlers.NewReviewRepository(db.DB)
	r.Handle("POST", "/api/reviews", middlewares.IsAuthorized(), handlers.SubmitReviewHandler(reviewRepo))


	err = r.Run()
	if err != nil {
		return
	}
	fmt.Println(db.DB.Config)

}
