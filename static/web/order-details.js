let currentRating = 0;

function setRating(rating) {
    currentRating = rating;
    document.getElementById('rating').value = rating; // Сохраняем значение в input
    updateStars();
}

function updateStars() {
    const stars = document.querySelectorAll('.star');
    stars.forEach(star => {
        const starValue = parseInt(star.getAttribute('data-value'));
        if (starValue <= currentRating) {
            star.classList.add('filled');
        } else {
            star.classList.remove('filled');
        }
    });
}

// Функция для открытия модального окна
function openReviewModal(itemID, itemName) {
    document.getElementById('review-item-id').value = itemID;
    document.getElementById('review-item-name').innerText = `Leave a Review for "${itemName}"`;
    document.getElementById('review-modal').style.display = "block";
}

// Функция для закрытия модального окна
function closeReviewModal() {
    document.getElementById('review-modal').style.display = "none";
    currentRating = 0; // Сброс звёзд
    updateStars();
}

// Функция для отправки формы
function submitReview(event) {
    event.preventDefault();

    const itemId = document.getElementById("review-item-id").value;
    const rating = document.getElementById("rating").value;
    const comment = document.getElementById("comment").value;

    const reviewData = {
        item_id: parseInt(itemId),
        rating: parseInt(rating),
        comment: comment.trim(),
    };

    console.log("Отправляемые данные:", reviewData);

    fetch("/api/reviews", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(reviewData),
    })
    .then((response) => response.json())
    .then((data) => {
        console.log("Ответ сервера:", data);
        if (data.success) {
            alert("Review submitted successfully!");
            closeReviewModal();
            location.reload();
        } else {
            alert("Error submitting review: " + data.error);
        }
    })
    .catch((error) => {
        console.error("Ошибка запроса:", error);
        alert("An unexpected error occurred.");
    });
}


