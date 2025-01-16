document.addEventListener("DOMContentLoaded", function() {

    const productCards = document.querySelectorAll(".product-card");

    productCards.forEach(card => {
        card.addEventListener("click", function(event) {
            const itemID = this.getAttribute("data-item-id");
            window.location.href = `/item/${itemID}`;
        });
    });

    // Находим все кнопки "Add to Cart"
    const cartButtons = document.querySelectorAll(".btn-add-to-cart");
    const notification = document.getElementById("notification");
    const notificationMessage = document.getElementById("notification-message");

    cartButtons.forEach(button => {
        button.addEventListener("click", function() {
        event.stopPropagation(); // Остановить всплытие события, чтобы не открывалась страница товара
            const itemID = this.getAttribute("data-item-id");
            const quantity = 1; // Добавляем по одному товару за раз

            // Проверяем, в корзине ли уже товар
            if (this.classList.contains("in-cart")) {
                // Перенаправляем на страницу корзины
                window.location.href = "/cart";
                return;
            }

            // Отправляем запрос на сервер для добавления товара в корзину
            fetch("/cart/add", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    item_id: parseInt(itemID),
                    quantity: quantity
                })
            })
            .then(response => response.json())
            .then(data => {
                if (data.message === "Item added to cart") {
                    button.classList.add("in-cart");
                        button.innerHTML = "In Cart";
                    showNotification("Item added to cart");
                } else {
                    alert("Error adding item to cart.");
                }
            })
            .catch(error => {
                console.error("Error:", error);
                alert("An error occurred. Please try again.");
            });
        });
    });

    const likeButtons = document.querySelectorAll(".btn-like");

    likeButtons.forEach(likeButton => {
    likeButton.addEventListener("click", function (event) {
        event.stopPropagation(); // Остановить всплытие события, чтобы не открывалась страница товара
        const itemID = this.getAttribute("data-item-id");
        const likedIcon = this.querySelector(".liked-icon1");
        // Отправляем запрос на сервер
        fetch("/liked/add", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ item_id: parseInt(itemID) }),
        })
            .then(response => response.json())
            .then(data => {
                console.log("Response from server:", data);

                // Меняем цвет сердечка локально
                if (likedIcon.classList.contains("liked")) {
                    likedIcon.classList.remove("liked");
                    showNotification("Removed from favorites");
                } else {
                    likedIcon.classList.add("liked");
                    showNotification("Added to favorites");
                }
                 // Показать всплывающее уведомление о добавлении товара в избранное
            })
            .catch(error => error("Fetch error:", error));
    });
});

// Функция для отображения уведомления
function showNotification(message) {
    notificationMessage.textContent = message;
    notification.classList.add("visible");
    notification.classList.remove("hidden");

    // Скрыть уведомление через 3 секунды
    setTimeout(() => {
        notification.classList.add("hidden");
        notification.classList.remove("visible");
    }, 2000);
}

    

    
});

function updateCartQuantity(itemID, quantity) {
    fetch("/cart/add", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ item_id: parseInt(itemID), quantity: quantity })
    }).then(response => response.json()).then(data => {
        if (data.message !== "Item added to cart") {
            console.error("Ошибка обновления количества в корзине:", data.error);
        }
    }).catch(error => console.error("Ошибка:", error));
}

