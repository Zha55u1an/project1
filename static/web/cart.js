document.addEventListener("DOMContentLoaded", function() {
    // Код для удаления из корзины
    const recentlyViewedContainer = document.querySelector('.recently-viewed-items');

    recentlyViewedContainer.addEventListener('wheel', (e) => {
        e.preventDefault();
        recentlyViewedContainer.scrollLeft += e.deltaY;
    });
    
});

function removeFromCart(itemId) {
    fetch('/cart/remove', {
        method: 'DELETE', // Или 'POST', если вы используете метод POST для удаления
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ item_id: itemId }),
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        location.reload(); // Обновить страницу, чтобы отобразить изменения в корзине
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

// Увеличивает или уменьшает количество на 1
function updateQuantity(itemId, change) {
    const quantityInput = document.getElementById(`quantity-${itemId}`);
    let currentQuantity = parseInt(quantityInput.value);

    // Обновляем количество с учетом границы минимума
    currentQuantity = Math.max(1, currentQuantity + change);
    quantityInput.value = currentQuantity;

    // Вызываем функцию для обработки изменений количества
    updateCart(itemId, currentQuantity);
}

// Устанавливает количество на основе ввода пользователя
function setQuantity(itemId, quantity) {
    const quantityInput = document.getElementById(`quantity-${itemId}`);
    let currentQuantity = parseInt(quantity);

    // Проверка на минимальное значение 1
    if (isNaN(currentQuantity) || currentQuantity < 1) {
        currentQuantity = 1;
    }
    quantityInput.value = currentQuantity;

    // Обновляем корзину на сервере
    updateCart(itemId, currentQuantity);
}

// Функция для обновления количества на сервере (асинхронно)
async function updateCart(itemId, quantity) {
    try {
        const response = await fetch('/update-cart-quantity', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ item_id: itemId, quantity: quantity })
        });

        if (!response.ok) {
            throw new Error('Ошибка при обновлении количества');
        }

        const data = await response.json();
        console.log(data.message);

        // Обновляем итоговую сумму корзины на странице
        if (data.total_amount !== undefined) {
            document.querySelector('.summary-total').innerText = `Total: ${data.total_amount}₸`;
            document.querySelector('.summary-detail').innerText = `Subtotal: ${data.total_amount}₸`;
        }
    } catch (error) {
        console.error(error);
        alert('Ошибка при обновлении корзины');
    }
}

function scrollProducts(direction) {
    const productGrid = document.querySelector('.product-grid');
    const scrollAmount = 300; // Количество пикселей для прокрутки
    productGrid.scrollBy({ left: direction * scrollAmount, behavior: 'smooth' });
}


// Код для работы с выбором адреса доставки
let selectedPoint = null;
let map; // Глобальная переменная для карты
let mapInitialized = false; // Флаг для проверки, была ли карта уже инициализирована


function openDeliveryModal() {
    document.getElementById('delivery-modal').style.display = 'flex';
    if (!mapInitialized) {
        ymaps.ready(initMap); // Создаём карту только один раз
        mapInitialized = true;
    }
}

function closeDeliveryModal() {
    document.getElementById('delivery-modal').style.display = 'none';
}




