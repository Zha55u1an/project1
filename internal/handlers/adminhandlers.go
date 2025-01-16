package handlers

import (
	"go_project/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

// NewAdminRepository создаёт новый экземпляр AdminRepository
func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetDashboardStats() (models.AdminDashboardStats, error) {
	var stats models.AdminDashboardStats

	// Считаем количество пользователей
	if err := r.db.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return stats, err
	}

	// Считаем количество товаров
	if err := r.db.Model(&models.Item{}).Count(&stats.TotalProducts).Error; err != nil {
		return stats, err
	}

	// Считаем общее количество заказов
	if err := r.db.Model(&models.Order{}).Count(&stats.TotalOrders).Error; err != nil {
		return stats, err
	}

	// Считаем заказы по статусам
	if err := r.db.Model(&models.Order{}).Where("status = ?", "pending").Count(&stats.PendingOrders).Error; err != nil {
		return stats, err
	}
	if err := r.db.Model(&models.Order{}).Where("status = ?", "completed").Count(&stats.CompletedOrders).Error; err != nil {
		return stats, err
	}

	return stats, nil
}


func AdminDashboardHandler(repo *AdminRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем статистику через репозиторий
		stats, err := repo.GetDashboardStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load dashboard stats"})
			return
		}

		// Передаём статистику в шаблон
		c.HTML(http.StatusOK, "admin-dashboard.html", gin.H{
			"TotalOrders":     stats.TotalOrders,
			"TotalProducts":   stats.TotalProducts,
			"TotalUsers":      stats.TotalUsers,
			"PendingOrders":   stats.PendingOrders,
			"CompletedOrders": stats.CompletedOrders,
		})
	}
}

func (r *AdminRepository) GetUsersByRole(role string) ([]models.User, error) {
	var users []models.User
	if err := r.db.Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func AdminUsersHandler(repo *AdminRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем пользователей с ролью "buyer"
		buyers, err := repo.GetUsersByRole("buyer")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load buyers"})
			return
		}

		// Получаем пользователей с ролью "seller"
		sellers, err := repo.GetUsersByRole("seller")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load sellers"})
			return
		}

		// Передаем данные в шаблон
		c.HTML(http.StatusOK, "admin-users.html", gin.H{
			"Buyers":  buyers,
			"Sellers": sellers,
		})
	}
}

func AdminOrdersHandler(repo *AdminRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем все заказы из базы данных
		var orders []models.Order

		// Загружаем все заказы, в том числе с их статусами и заказанными товарами
		if err := repo.db.Preload("OrderItems.Item").Find(&orders).Error; err != nil {
			log.Printf("Error fetching orders: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load orders"})
			return
		}

		// Передаем данные о заказах в шаблон
		c.HTML(http.StatusOK, "admin-order.html", gin.H{
			"Orders": orders, // Передаем список всех заказов
		})
	}
}

func (r *AdminRepository) GetOrderDetails(orderID uint) (models.Order, error) {
	var order models.Order
	// Загружаем заказ с товарами и их данными
	if err := r.db.Preload("OrderItems.Item").First(&order, orderID).Error; err != nil {
		return order, err
	}
	return order, nil
}
func AdminOrderDetailsHandler(repo *AdminRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID заказа из URL
		orderIDStr := c.Param("id")
		orderID, err := strconv.ParseUint(orderIDStr, 10, 32) // Преобразуем строку в uint
		if err != nil {
			log.Printf("Invalid order ID: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		// Получаем заказ с деталями
		order, err := repo.GetOrderDetails(uint(orderID))
		if err != nil {
			log.Printf("Order not found: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		// Форматируем дату
		formattedDate := order.CreatedAt.Format("2006-01-02") // Дата в формате YYYY-MM-DD
		formattedTime := order.CreatedAt.Format("15:04") // Время в формате HH:MM

		// Передаем данные в шаблон
		c.HTML(http.StatusOK, "admin-order-details.html", gin.H{
			"Order": gin.H{
				"ID":        order.ID,
				"Date":      formattedDate,
				"Time":      formattedTime,
				"Status":    order.Status,
				"OrderItems": order.OrderItems, // Передаем товары
			},
		})
	}
}

func (r *AdminRepository) UpdateOrderStatus(orderID uint, status string) error {
	var order models.Order
	// Находим заказ по ID
	if err := r.db.First(&order, orderID).Error; err != nil {
		return err
	}

	// Обновляем статус
	order.Status = status
	if err := r.db.Save(&order).Error; err != nil {
		return err
	}

	return nil
}
func AdminOrderUpdateHandler(repo *AdminRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID заказа из URL
		orderIDStr := c.Param("id")
		orderID, err := strconv.ParseUint(orderIDStr, 10, 32) // Преобразуем строку в uint
		if err != nil {
			log.Printf("Invalid order ID: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		// Получаем заказ с деталями
		order, err := repo.GetOrderDetails(uint(orderID))
		if err != nil {
			log.Printf("Order not found: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		// Обрабатываем POST запрос на изменение статуса заказа
		if c.Request.Method == "POST" {
			// Получаем новый статус из формы
			newStatus := c.PostForm("status")
			if newStatus == "" {
				log.Printf("Status is required")
				c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
				return
			}

			// Обновляем статус заказа
			err := repo.UpdateOrderStatus(uint(orderID), newStatus)
			if err != nil {
				log.Printf("Error updating order status: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating order status"})
				return
			}

			// Перенаправляем на страницу с деталями обновленного заказа
			c.Redirect(http.StatusSeeOther, "/admin/order-details/"+orderIDStr)
			return
		}

		// Форматируем дату
		formattedDate := order.CreatedAt.Format("2006-01-02") // Дата в формате YYYY-MM-DD
		formattedTime := order.CreatedAt.Format("15:04") // Время в формате HH:MM

		// Передаем данные в шаблон для отображения
		c.HTML(http.StatusOK, "admin-order-update.html", gin.H{
			"Order": gin.H{
				"ID":        order.ID,
				"Date":      formattedDate,
				"Time":      formattedTime,
				"Status":    order.Status,
				"OrderItems": order.OrderItems, // Передаем товары
			},
		})
	}
}



