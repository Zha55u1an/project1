package handlers

import (
	"fmt"
	"go_project/internal/models"
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
// CreateOrder handles POST request to create a new order
func (repo *OrderRepository) CreateOrder(c *gin.Context) {
	var orderItems []models.OrderItem

	if err := c.BindJSON(&orderItems); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	order := models.Order{
		UserID: userID.(uint),
		// CreatedAt:  time.Now(),
		Status: "created",
	}
	result := repo.db.Create(&order)
	if result.Error != nil {
		// Выводим ошибку в консоль
		fmt.Println("Error creating order:", result.Error.Error())
		// Возвращаем ошибку в ответе сервера
		c.JSON(http.StatusInternalServerError, gin.H{"error LLL": result.Error.Error()})
		return
	}
	// Создаем заказы на основе полученных данных
	for _, item := range orderItems {
		// Создаем новый заказ для каждого элемента

		// Создаем заказ в базе данных

		// Получаем orderID после создания заказа
		orderID := order.ID

		// Привязываем каждый OrderItem к orderID
		item.OrderID = orderID
		result = repo.db.Create(&item)
		if result.Error != nil {
			// Выводим ошибку в консоль
			fmt.Println("Error creating order item:", result.Error.Error())
			// Возвращаем ошибку в ответе сервера
			c.JSON(http.StatusInternalServerError, gin.H{"error MMM": result.Error.Error()})
			return
		}

		// Обновляем статус заказа на "created"
		repo.db.Model(&order).Update("status", "created")
	}

	// Возвращаем ответ с информацией о заказах
	c.JSON(http.StatusCreated, gin.H{
		"message": "Orders created successfully",
	})
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
