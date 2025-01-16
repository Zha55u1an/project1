package handlers

import (
	"bytes"
	"go_project/internal/models"
	"go_project/pkg/db"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// GetAllOrders handles GET request to fetch all orders
func (repo *OrderRepository) GetAllOrders(c *gin.Context) {
	var orders []models.Order
	result := repo.db.Preload("OrderItems").Find(&orders)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// GetOrderByID handles GET request to fetch an order by ID
func (repo *OrderRepository) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	result := repo.db.Preload("OrderItems").First(&order, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// CreateOrder handles POST request to create a new order
// Создание заказа с товарами из корзины
func CreateOrder(repo *OrderRepository, cartRepo *CartRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Starting CreateOrderHandler")

		// Получаем userID из контекста
		userID, exists := c.Get("userID")
		if !exists {
			log.Println("User ID not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		userIDUint, ok := userID.(uint)
		if !ok {
			log.Printf("Invalid user ID type: %v", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			return
		}
		log.Printf("User ID: %d", userIDUint)

		// Получаем товары из корзины
		cartItems, err := cartRepo.GetCartItemsByUserID(userIDUint)
		if err != nil {
			log.Printf("Error fetching cart items: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch cart items"})
			return
		}

		if len(cartItems) == 0 {
			log.Println("Cart is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Your cart is empty"})
			return
		}

		// Создаём заказ с товарами из корзины
		orderItems := make([]models.OrderItem, len(cartItems))
		for i, cartItem := range cartItems {
			orderItems[i] = models.OrderItem{
				ItemID:   cartItem.ItemID,
				Quantity: cartItem.Quantity,
				Price:    cartItem.Item.Price,
			}
		}

		order := models.Order{
			UserID:     userIDUint,
			Status:     "created",
			OrderItems: orderItems, // Добавляем товары в заказ
		}

		if err := repo.db.Create(&order).Error; err != nil {
			log.Printf("Error creating order with items: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating order"})
			return
		}

		// Очищаем корзину
		if err := cartRepo.ClearCartByUserID(userIDUint); err != nil {
			log.Printf("Error clearing cart: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to clear cart"})
			return
		}

		// Перенаправляем пользователя на страницу заказов
		c.Redirect(http.StatusSeeOther, "/my-orders")
	}
}

// Получить товары из корзины пользователя
func (r *CartRepository) GetCartItemsByUserID(userID uint) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	err := r.db.Preload("Item").Where("user_id = ?", userID).Find(&cartItems).Error
	return cartItems, err
}

// Очистить корзину пользователя
func (r *CartRepository) ClearCartByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Unscoped().Delete(&models.CartItem{}).Error
}

func MyOrdersHandler(repo *OrderRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		userIDUint := userID.(uint)

		var orders []models.Order
		if err := repo.db.Preload("OrderItems").Where("user_id = ?", userIDUint).Find(&orders).Error; err != nil {
			log.Printf("Error fetching orders: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch orders"})
			return
		}

		// Рассчитываем сумму заказа
		type OrderView struct {
			ID        uint
			CreatedAt string
			Status    string
			Total     float64
		}

		var orderViews []OrderView
		for _, order := range orders {
			total := 0.0
			for _, item := range order.OrderItems {
				total += item.Price * float64(item.Quantity)
			}

			orderViews = append(orderViews, OrderView{
				ID:        order.ID,
				CreatedAt: order.CreatedAt.Format("02 Jan 2006 15:04"), // Форматируем дату и время
				Status:    order.Status,
				Total:     total,
			})
		}

		// Передаём данные в шаблон
		c.HTML(http.StatusOK, "my-orders.html", gin.H{
			"Orders": orderViews,
		})
	}
}


// UpdateOrder handles PUT request to update an existing order
func (repo *OrderRepository) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := repo.db.Model(&models.Order{}).Where("id = ?", id).Updates(&order)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	repo.db.Where("id=?", id).First(&order)
	c.JSON(http.StatusOK, order)
}

// DeleteOrder handles DELETE request to delete an order by ID
func (repo *OrderRepository) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	result := repo.db.Delete(&models.Order{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}


func GetPickupPoints(c *gin.Context) {
	// Используем глобальную переменную db.DB для работы с базой данных
	var points []models.PickupPoint

	// Выполняем запрос к базе данных для получения пунктов выдачи
	result := db.DB.Find(&points)

	// Обрабатываем возможную ошибку
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка выполнения запроса: " + result.Error.Error()})
		return
	}

	// Проверяем, если записи не найдены
	if len(points) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Пункты выдачи не найдены"})
		return
	}

	// Возвращаем пункты выдачи как JSON
	c.JSON(http.StatusOK, points)
}


func GetOrderDetailsHandler(repo *OrderRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID заказа из URL
		orderID := c.Param("id")

		// Находим заказ с деталями
		var order models.Order
		// Загружаем заказ с товарами и их данными
        if err := repo.db.Preload("OrderItems.Item").First(&order, orderID).Error; err != nil {
            log.Printf("Order not found: %v", err)
            c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
            return
        }

        // Загружаем изображения для каждого товара в заказе
        for i, orderItem := range order.OrderItems {
            var itemImage models.ItemImage
            // Находим изображение для каждого товара
            if err := repo.db.Where("item_id = ?", orderItem.ItemID).First(&itemImage).Error; err == nil {
                // Добавляем изображение в товар
                order.OrderItems[i].Item.ImagePath = itemImage.ImagePath
            }
        }

        // Форматируем дату и время
        formattedTime := order.CreatedAt.Format("15:04")       // Часы и минуты (21:57)
        formattedDate := order.CreatedAt.Format("2006-01-02") // Дата (2024-12-15)

        // Передаём данные в шаблон
        c.HTML(http.StatusOK, "order-details.html", gin.H{
            "Order": gin.H{
                "ID":     order.ID,
                "Time":   formattedTime,
                "Date":   formattedDate,
                "Status": order.Status,
                "OrderItems": order.OrderItems, // Передаём товары
            },
        })
	}
}


// ReviewRepository структура для работы с отзывами
type ReviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository создает новый экземпляр ReviewRepository
func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) CreateReview(review *models.Review) error {
    if err := r.db.Create(review).Error; err != nil {
        log.Println("Error creating review:", err)
        return err
    }
    return nil
}


// SubmitReviewHandler обрабатывает добавление отзыва
func SubmitReviewHandler(repo *ReviewRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Чтение и логирование тела запроса
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("Error reading body:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read body"})
			return
		}
		log.Println("Raw Request Body:", string(body))

		// Восстанавливаем тело запроса
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Парсинг JSON
		var request struct {
			ItemID  uint   `json:"item_id"`
			Rating  int    `json:"rating"`
			Comment string `json:"comment"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Получение userID из middleware
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Создание отзыва
		review := models.Review{
			UserID:  userID.(uint),
			ItemID:  request.ItemID,
			Rating:  request.Rating,
			Comment: request.Comment,
		}

		// Сохранение через репозиторий
		if err := repo.CreateReview(&review); err != nil {
			log.Println("Error saving review:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save review"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Review submitted successfully"})
	}
}





