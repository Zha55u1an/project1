package handlers

import (
	"errors"
	"go_project/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

// Обработчик для добавления товара в корзину
func AddToCartHandler(repo *CartRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			ItemID   uint `json:"item_id"`
			Quantity int  `json:"quantity"`
		}

		// Проверяем JSON запрос
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Проверка, что ItemID и Quantity корректны
        if request.ItemID == 0 || request.Quantity <= 0 {
            log.Println("Invalid item_id or quantity in request:", request.ItemID, request.Quantity)
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item_id or quantity"})
            return
        }

		// Извлекаем ID пользователя из контекста (должен быть установлен аутентификационным middleware)
		userID, exists := c.Get("userID")
		if !exists {
			log.Println("User not authenticated")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Преобразуем userID в тип uint
		userIDUint, ok := userID.(uint)
		if !ok {
			log.Println("Invalid user ID:", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			return
		}

		// Добавляем товар в корзину
		if err := repo.AddToCart(userIDUint, request.ItemID, request.Quantity); err != nil {
			log.Println("Unable to add item to cart:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add item to cart"})
			return
		}

		log.Println("Item added to cart successfully:", request.ItemID)
		c.JSON(http.StatusOK, gin.H{"message": "Item added to cart"})
	}
}

// Функция для добавления товара в корзину в репозитории
func (repo *CartRepository) AddToCart(userID, itemID uint, quantity int) error {
	
	var cartItem models.CartItem
	result := repo.db.Where("user_id = ? AND item_id = ?", userID, itemID).First(&cartItem)

	if result.Error == nil {
		// Если товар уже в корзине, увеличиваем количество
		cartItem.Quantity += quantity
		if err := repo.db.Save(&cartItem).Error; err != nil {
			log.Println("Error updating cart item:", err)
			return err
		}
		log.Println("Updated cart item quantity:", cartItem)
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Если товара нет в корзине, создаем новую запись
		newCartItem := models.CartItem{
			UserID:   userID,
			ItemID:   itemID,
			Quantity: quantity,
		}
		if err := repo.db.Create(&newCartItem).Error; err != nil {
			log.Println("Error creating new cart item:", err)
			return err
		}
		log.Println("Created new cart item:", newCartItem)
	} else {
		log.Println("Database error when checking for existing cart item:", result.Error)
		return result.Error
	}
	return nil
}

func (repo *CartRepository) GetCartItems(userID uint) ([]models.CartItem, error) {
    var cartItems []models.CartItem
    
    // Извлекаем товары из корзины с предзагрузкой связанных данных (товар и его изображение)
    if err := repo.db.Preload("Item").Preload("Item.Images").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
        log.Println("Error retrieving cart items:", err)
        return nil, err
    }
    
    // Заполняем путь к изображению для каждого товара
    for i := range cartItems {
        if len(cartItems[i].Item.Images) > 0 {
            // Предполагаем, что товары могут иметь несколько изображений
            cartItems[i].Item.ImagePath = cartItems[i].Item.Images[0].ImagePath // Устанавливаем первое изображение
        }
    }

    log.Println("Retrieved cart items for user:", userID, cartItems)
    return cartItems, nil
}


func ViewCartHandler(repo *ItemRepository, cartRepo *CartRepository, likedRepo *LikedRepository, recentlyViewedRepo *RecentlyViewedRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("userID")
        if !exists {
            log.Println("Пользователь не аутентифицирован")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
            return
        }

        userIDUint, ok := userID.(uint)
        if !ok {
            log.Println("Неверный ID пользователя:", userID)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный ID пользователя"})
            return
        }

        // Получаем товары в корзине
        cartItems, err := cartRepo.GetCartItems(userIDUint)
        if err != nil {
            log.Println("Не удалось получить товары в корзине:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить товары в корзине"})
            return
        }

        // Вычисление общей суммы
        totalAmount := 0.0
        for _, item := range cartItems {
            totalAmount += float64(item.Item.Price) * float64(item.Quantity)
        }

        // Получаем недавно просмотренные товары с изображениями
        recentlyViewedItems, err := recentlyViewedRepo.GetRecentlyViewed(userIDUint, 10) // Лимитируем 10 товарами
        if err != nil {
            log.Println("Не удалось получить недавно просмотренные товары:", err)
            recentlyViewedItems = []models.Item{} // Если ошибка, возвращаем пустой список
        }

        // Для каждого недавно просмотренного товара загружаем изображение, рейтинг и количество отзывов
        for i := range recentlyViewedItems {
            // Загрузка изображений для товара
            var itemImage models.ItemImage
            if err := recentlyViewedRepo.db.Where("item_id = ?", recentlyViewedItems[i].ID).First(&itemImage).Error; err == nil {
                recentlyViewedItems[i].ImagePath = itemImage.ImagePath
            }

            // Средний рейтинг для недавно просмотренного товара
            var averageRating float64
            recentlyViewedRepo.db.Model(&models.Review{}).Where("item_id = ?", recentlyViewedItems[i].ID).Select("AVG(rating)").Scan(&averageRating)
            recentlyViewedItems[i].AverageRating = averageRating

            // Количество отзывов для товара
            var reviewCount int64
            recentlyViewedRepo.db.Model(&models.Review{}).Where("item_id = ?", recentlyViewedItems[i].ID).Count(&reviewCount)
            recentlyViewedItems[i].ReviewCount = reviewCount
        }

        // Проверяем, какие товары в избранных у текущего пользователя
        likedItems := []models.LikedItem{}
        likedRepo.db.Where("user_id = ?", userIDUint).Find(&likedItems)

        // Создаем карту, где ключ — `item_id`, а значение — true, если товар в избранных
        likedMap := make(map[uint]bool)
        for _, likedItem := range likedItems {
            likedMap[likedItem.ItemID] = true
        }

        // Создаем карту товаров в корзине для быстрого доступа
        cartMap := make(map[uint]bool)
        for _, item := range cartItems {
            cartMap[item.ItemID] = true
        }

        // Подготовка данных для шаблона
        data := gin.H{
            "User":              userIDUint,
            "CartItems":         cartItems,
            "TotalAmount":       totalAmount,
            "RecentlyViewed":    recentlyViewedItems,
            "LikedMap":          likedMap,  // Передаем карту избранных
            "CartMap":           cartMap,   // Передаем карту товаров в корзине
        }

        log.Println("Рендеринг страницы корзины для пользователя:", userIDUint)
        c.HTML(http.StatusOK, "cart.html", data)
    }
}






func ViewCheckoutHandler(repo *CartRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("userID")
        if !exists {
            log.Println("Пользователь не аутентифицирован")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
            return
        }

        userIDUint, ok := userID.(uint)
        if !ok {
            log.Println("Неверный ID пользователя:", userID)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный ID пользователя"})
            return
        }

        // Получение товаров в корзине
        cartItems, err := repo.GetCartItems(userIDUint)
        if err != nil {
            log.Println("Не удалось получить товары в корзине:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить товары в корзине"})
            return
        }

        // Вычисление общей суммы
        totalAmount := 0.0
        for _, item := range cartItems {
            totalAmount += float64(item.Item.Price) * float64(item.Quantity)
        }

        log.Println("Рендеринг страницы оплаты для пользователя:", userIDUint)
        c.HTML(http.StatusOK, "checkout.html", gin.H{
            "TotalAmount": totalAmount, 
        })
    }
}


// Обработчик для удаления товара из корзины
func RemoveFromCartHandler(repo *CartRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        var request struct {
            ItemID uint `json:"item_id"`
        }

        // Проверяем JSON запрос
        if err := c.ShouldBindJSON(&request); err != nil {
            log.Println("Error binding JSON:", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
            return
        }

        // Извлекаем ID пользователя из контекста
        userID, exists := c.Get("userID")
        if !exists {
            log.Println("User not authenticated")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
            return
        }

        userIDUint, ok := userID.(uint)
        if !ok {
            log.Println("Invalid user ID:", userID)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
            return
        }

        // Удаляем товар из корзины
        if err := repo.RemoveFromCart(userIDUint, request.ItemID); err != nil {
            log.Println("Error removing item from cart:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove item from cart"})
            return
        }

        log.Println("Item removed from cart successfully:", request.ItemID)
        c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
    }
}

// Функция для удаления товара из корзины в репозитории
func (repo *CartRepository) RemoveFromCart(userID, itemID uint) error {
    // Здесь выполняем удаление товара из корзины
    if err := repo.db.Unscoped().Where("user_id = ? AND item_id = ?", userID, itemID).Delete(&models.CartItem{}).Error; err != nil {
        return err
    }
    return nil
}


// Обработчик для обновления количества товара в корзине
func UpdateCartQuantityHandler(repo *CartRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        var request struct {
            ItemID   uint `json:"item_id"`
            Quantity int  `json:"quantity"`
        }

        // Проверяем JSON запрос
        if err := c.ShouldBindJSON(&request); err != nil {
            log.Println("Ошибка привязки JSON:", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
            return
        }

        // Проверка, что ItemID и Quantity корректны
        if request.ItemID == 0 || request.Quantity <= 0 {
            log.Println("Неверный item_id или quantity в запросе:", request.ItemID, request.Quantity)
            c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный item_id или quantity"})
            return
        }

        // Извлекаем ID пользователя
        userID, exists := c.Get("userID")
        if !exists {
            log.Println("Пользователь не аутентифицирован")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
            return
        }

        userIDUint, ok := userID.(uint)
        if !ok {
            log.Println("Неверный ID пользователя:", userID)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный ID пользователя"})
            return
        }

        // Обновляем количество товара в корзине
        if err := repo.UpdateCartQuantity(userIDUint, request.ItemID, request.Quantity); err != nil {
            log.Println("Ошибка при обновлении количества товара в корзине:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить количество товара в корзине"})
            return
        }

        // Получаем обновленную сумму корзины
        totalAmount, err := repo.GetCartTotalAmount(userIDUint)
        if err != nil {
            log.Println("Ошибка при получении общей суммы корзины:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить общую сумму корзины"})
            return
        }

        log.Println("Количество товара в корзине обновлено успешно:", request.ItemID)
        c.JSON(http.StatusOK, gin.H{"message": "Количество товара обновлено", "total_amount": totalAmount})
    }
}


// Функция для обновления количества товара в корзине в репозитории
func (repo *CartRepository) UpdateCartQuantity(userID, itemID uint, quantity int) error {
    var cartItem models.CartItem
    result := repo.db.Where("user_id = ? AND item_id = ?", userID, itemID).First(&cartItem)

    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            log.Println("Товар не найден в корзине для обновления:", itemID)
            return errors.New("товар не найден в корзине")
        }
        log.Println("Ошибка базы данных при поиске товара в корзине:", result.Error)
        return result.Error
    }

    // Обновляем количество товара
    cartItem.Quantity = quantity
    if err := repo.db.Save(&cartItem).Error; err != nil {
        log.Println("Ошибка при сохранении обновленного количества товара:", err)
        return err
    }

    log.Println("Обновленное количество товара в корзине:", cartItem)
    return nil
}

func (repo *CartRepository) GetCartTotalAmount(userID uint) (float64, error) {
    var totalAmount float64
    err := repo.db.Model(&models.CartItem{}).
        Select("SUM(cart_items.quantity * items.price)").
        Joins("JOIN items ON items.id = cart_items.item_id").
        Where("cart_items.user_id = ?", userID).
        Scan(&totalAmount).Error

    if err != nil {
        return 0, err
    }
    return totalAmount, nil
}

func GetCartItemCountHandler(repo *CartRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("userID")
        var cartItemsCount int

        if exists {
            userIDUint, ok := userID.(uint)
            if ok {
                cartItems, err := repo.GetCartItems(userIDUint)
                if err == nil {
                    cartItemsCount = len(cartItems)
                }
            }
        }

        c.HTML(http.StatusOK, "base.html", gin.H{
            "CartItemsCount": cartItemsCount, // Передаем количество товаров
        })
    }
}


func CartMiddleware(cartRepo *CartRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        var cartItemsCount int
        if userID, exists := c.Get("userID"); exists {
            userIDUint, ok := userID.(uint)
            if ok {
                cartItems, err := cartRepo.GetCartItems(userIDUint)
                if err == nil {
                    cartItemsCount = len(cartItems)
                }
            }
        }
        c.Set("CartItemsCount", cartItemsCount)
        c.Next()
    }
}


func (repo *RecentlyViewedRepository) GetRecentlyViewed(userID uint, limit int) ([]models.Item, error) {
    var recentlyViewed []models.RecentlyViewedItem
    var items []models.Item

    // Получаем последние просмотренные товары
    if err := repo.db.
        Where("user_id = ?", userID).
        Order("view_time DESC").
        Limit(limit).
        Find(&recentlyViewed).Error; err != nil {
        return nil, err
    }

    // Извлекаем информацию о товарах
    for _, viewed := range recentlyViewed {
        var item models.Item
        if err := repo.db.First(&item, viewed.ItemID).Error; err == nil {
            items = append(items, item)
        }
    }

    return items, nil
}