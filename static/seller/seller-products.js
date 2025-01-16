document.getElementById('addProductForm').addEventListener('submit', async function(e) {
    e.preventDefault(); // Останавливаем перезагрузку страницы

    // Получаем данные из формы
    const productName = document.getElementById('productName').value;
    const productPrice = document.getElementById('productPrice').value;
    const productCategory = document.getElementById('productCategory').value;
    const productDescription = document.getElementById('productDescription').value;
    const productImage = document.getElementById('productImage').files[0];

    // Создаем FormData для отправки файла и данных
    const formData = new FormData();
    formData.append('productName', productName);
    formData.append('productPrice', productPrice);
    formData.append('productCategory', productCategory);
    formData.append('productDescription', productDescription);
    if (productImage) {
        formData.append('productImage', productImage);
    }

    try {
        const response = await fetch('/seller/products', {
            method: 'POST',
            body: formData
        });

        const result = await response.json();

        if (response.ok) {
            alert('Product added successfully!');
            // Обновляем список товаров
            loadProducts();
        } else {
            alert('Error: ' + result.error);
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Something went wrong!');
    }
});

document.addEventListener('DOMContentLoaded', function() {
    console.log("DOMContentLoaded: Скрипт загружен.");

    // Использование делегирования для обработки кликов по кнопкам с классом delete-button
    document.querySelector('table').addEventListener('click', function(e) {
        console.log("Клик по таблице: ", e.target);
        
        if (e.target.closest('.delete-button')) {
            const itemId = e.target.closest('.delete-button').getAttribute('data-item-id');
            console.log("Найден товар для удаления. ID: ", itemId);

            const confirmation = confirm('Are you sure you want to delete this product?');
            console.log("Подтверждение удаления: ", confirmation);

            if (confirmation) {
                console.log("Запрос на удаление товара ID:", itemId);

                // Выполняем запрос на удаление товара
                fetch(`/seller/products/delete/${itemId}`, {
                    method: 'DELETE',
                })
                .then(response => {
                    if (response.ok) {
                        console.log("Товар успешно удален.");
                        alert('Product deleted successfully!');
                        loadProducts();  // Обновляем список товаров после удаления
                    } else {
                        console.error("Ошибка удаления товара, код ответа: ", response.status);
                        alert('Error deleting product.');
                    }
                })
                .catch(error => {
                    console.error("Ошибка запроса на удаление товара: ", error);
                    alert('There was an error deleting the product.');
                });
            } else {
                console.log("Удаление отменено пользователем.");
            }
        }
    });

    // Функция для загрузки списка товаров
    async function loadProducts() {
        console.log("Загружаем список товаров...");

        try {
            const response = await fetch('/seller/products', {
                method: 'GET'
            });
            console.log("Ответ от сервера с товарами: ", response);

            const products = await response.json();
            console.log("Продукты получены: ", products);

            const productTableBody = document.querySelector('#productTable tbody');
            productTableBody.innerHTML = ''; // Очищаем таблицу

            products.forEach(product => {
                console.log("Добавление товара в таблицу: ", product);

                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${product.name}</td>
                    <td>${product.price}₸</td>
                    <td>${product.category}</td>
                    <td>${product.description}</td>
                    <td><img src="${product.imagePath}" alt="${product.name}" width="50"></td>
                    <td>
                        <button class="btn-edit" onclick="editProduct(${product.id})">Edit</button>
                        <button class="delete-button" data-item-id="${product.id}">Delete</button>
                    </td>
                `;
                productTableBody.appendChild(row);
            });
        } catch (error) {
            console.error('Ошибка при загрузке продуктов: ', error);
        }
    }

    // Загружаем список товаров при загрузке страницы
    loadProducts();
});
