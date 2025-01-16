document.addEventListener("DOMContentLoaded", function () {
    // Получаем элемент для отображения сообщения
    const messageElement = document.getElementById("profile-message");
    const form = document.getElementById("profile-form");

    // Функция для отображения сообщений
    function showMessage(message, type) {
        messageElement.textContent = message;
        messageElement.className = type; // Для изменения стиля в зависимости от типа сообщения
    }

    // Функция для обработки ответа в JSON
    function handleJsonResponse(response) {
        if (!response.ok) {
            // Если статус ответа не ок, выводим ошибку
            return response.text().then(text => {
                showMessage("Произошла ошибка на сервере: " + text, "error");
                throw new Error("Server Error");
            });
        }
        return response.json();
    }

    // Загружаем данные профиля
    fetch("/profile")
    .then(handleJsonResponse)
    .then(data => {
        if (data.message) {
            // Если сервер вернул сообщение, показываем его
            showMessage(data.message, "error");
        } else {
            // Здесь можно добавить логику для отображения информации профиля
            // Пример:
            document.getElementById("surname").value = data.userInfo.Surname || "";
            document.getElementById("email").value = data.userInfo.Email || "";
        }
    })
    .catch(error => {
        console.error("Error:", error);
        showMessage("", "error");
    });

    // Обработчик отправки формы
    form.addEventListener("submit", function (event) {
        event.preventDefault(); // Отменяем стандартное поведение формы

        // Получаем данные из формы
        const formData = new FormData(form);
        const data = {};

        formData.forEach((value, key) => {
            data[key] = value;
        });

        // Отправляем данные на сервер
        fetch("/profile-update", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
        })
        .then(handleJsonResponse)
        .then(data => {
            if (data.message) {

                // Если сервер вернул сообщение, показываем его
                showMessage(data.message, "success");
                document.getElementById("profile-message").style.display = "flex";
            } else {
                // В случае успешного обновления, обновляем страницу
                window.location.reload();
            }
        })
        .catch(error => {
            console.error("Error:", error);
            showMessage("", "error");
        });
    });
});
