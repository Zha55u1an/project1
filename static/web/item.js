function updateMainImage(imagePath) {
    const mainImage = document.getElementById('main-image');
    mainImage.src = imagePath;
}

// Функция для создания звезд
function displayStars() {
    const ratings = document.querySelectorAll('.rating-stars');
    
    ratings.forEach((ratingElement) => {
        const rating = parseFloat(ratingElement.getAttribute('data-rating'));
        let starsHtml = '';
        
        // Создаем звезды на основе рейтинга
        for (let i = 1; i <= 5; i++) {
            if (i <= rating) {
                starsHtml += '<span class="star filled">★</span>'; // Заполненная звезда
            } else {
                starsHtml += '<span class="star">☆</span>'; // Пустая звезда
            }
        }

        // Вставляем звезды в элемент
        ratingElement.innerHTML = starsHtml;
    });
}

// Запуск функции после загрузки страницы
window.addEventListener('DOMContentLoaded', displayStars);
