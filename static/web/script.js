document.querySelector('.login-button').addEventListener('mouseenter', function () {
    document.querySelector('.dropdown-menu').style.display = 'block';
});

document.querySelector('.dropdown-menu').addEventListener('mouseleave', function () {
    document.querySelector('.dropdown-menu').style.display = 'none';
});

document.querySelector('.dropdown-menu').addEventListener('mouseenter', function () {
    document.querySelector('.dropdown-menu').style.display = 'block';
});

document.querySelector('.login-button').addEventListener('mouseleave', function () {
    document.querySelector('.dropdown-menu').style.display = 'none';
});

document.getElementById("logout-button").addEventListener("click", function () {
    fetch("/logout", {
        method: "POST",
        credentials: "include"
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            window.location.href = data.redirect; // Перенаправление
        } else {
            console.error("Logout failed:", data.message);
        }
    })
    .catch(error => {
        console.error("Error during logout:", error);
    });
});

function toggleMenu() {
    const sideMenu = document.getElementById("sideMenu");
    const overlay = document.getElementById("overlay");
    const hamburger = document.querySelector(".hamburger-menu");

    sideMenu.classList.toggle("open");
    overlay.classList.toggle("open");
    hamburger.classList.toggle("active");
}


function closeMenu() {
    const sideMenu = document.getElementById("sideMenu");
    const overlay = document.getElementById("overlay");
    const hamburger = document.querySelector(".hamburger-menu");

    sideMenu.classList.remove("open");
    overlay.classList.remove("open");
    hamburger.classList.remove("active");
}

document.addEventListener("DOMContentLoaded", () => {
    // Выполняем GET-запрос к API, чтобы получить адрес доставки
    fetch('/api/get-delivery-address', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to load delivery address');
        }
        return response.json();
    })
    .then(data => {
        // Формируем текст адреса
        const addressText = `${data.address_line}`;

        // Обновляем элемент с id "delivery-address" в хедере
        document.getElementById('delivery-address').textContent = `${addressText}`;
        
        // Показать ссылку для редактирования адреса
    })
    .catch(error => {
        console.error('Error fetching delivery address:', error);

        // Если не удалось загрузить адрес, показываем стандартное сообщение
        document.getElementById('delivery-address').style.display = 'none';
        document.getElementById('delivery-icon').style.display = 'none';
    });
});



