function showNotification(message) {
    const notification = document.getElementById("notification");
    const notificationMessage = document.getElementById("notification-message");
    notificationMessage.textContent = message;
    notification.classList.add("visible");
    notification.classList.remove("hidden");

    // Скрыть уведомление через 3 секунды
    setTimeout(() => {
        notification.classList.add("hidden");
        notification.classList.remove("visible");
    }, 2000);
}

function addToFavorites(itemID) {
    // Проверяем, что itemID существует и является допустимым значением
    if (!itemID || itemID === "null") {
        console.error("Invalid item ID:", itemID);
        return; // Прерываем выполнение, если itemID некорректный
    }

    // Отправляем запрос на сервер для добавления товара в избранное
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
                location.reload(); // Обновить страницу, чтобы отобразить изменения в корзине

            } else {
                likedIcon.classList.add("liked");
                showNotification("Added to favorites");
            }
             // Показать всплывающее уведомление о добавлении товара в избранное
        })
        .catch(error => error("Fetch error:", error));
}

