document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');
    
    if (loginForm) {
        loginForm.addEventListener('submit', async function(e) {
            e.preventDefault(); // Останавливаем стандартное поведение формы

            const email = loginForm.querySelector('input[name="email"]').value; // Получаем значение поля email
            const password = loginForm.querySelector('input[name="password"]').value; // Получаем значение поля password

            const userData = {
                email: email,
                password: password
            };

            try {
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(userData) // Отправляем данные как JSON
                });

                const result = await response.json(); // Парсим JSON ответ

                if (response.ok) {
                    // Перенаправляем на нужную страницу в зависимости от роли
                    window.location.href = result.redirect; // Редирект
                } else {
                    document.getElementById('loginMessage').textContent = result.error; // Выводим сообщение об ошибке
                    document.getElementById('loginMessage').style.color = 'red'; // Устанавливаем цвет ошибки
                }
            } catch (error) {
                console.error('Error:', error); // Логируем ошибку в консоль
                document.getElementById('loginMessage').textContent = 'Something went wrong!'; // Общая ошибка
                document.getElementById('loginMessage').style.color = 'red'; // Цвет ошибки
            }
        });
    } else {
        console.error('Login form not found!'); // Если форма не найдена
    }
});
// Функция для переключения видимости пароля
function togglePasswordVisibility(inputId, iconId) {
    const input = document.getElementById(inputId);
    const icon = document.getElementById(iconId).querySelector('svg');

    if (input.type === 'password') {
        input.type = 'text';
        icon.innerHTML = `
            <path strokeLinecap="round" strokeLinejoin="round" d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z" />
            <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
        `;
    } else {
        input.type = 'password';
        icon.innerHTML = `
            <path strokeLinecap="round" strokeLinejoin="round" d="M3.98 8.223A10.477 10.477 0 0 0 1.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.451 10.451 0 0 1 12 4.5c4.756 0 8.773 3.162 10.065 7.498a10.522 10.522 0 0 1-4.293 5.774M6.228 6.228 3 3m3.228 3.228 3.65 3.65m7.894 7.894L21 21m-3.228-3.228-3.65-3.65m0 0a3 3 0 1 0-4.243-4.243m4.242 4.242L9.88 9.88" />
        `;
    }
}

// Обработчики для полей пароля
document.getElementById('togglePassword').addEventListener('click', () => {
    togglePasswordVisibility('password', 'togglePassword');
});

document.getElementById('toggleConfirmPassword').addEventListener('click', () => {
    togglePasswordVisibility('confirmPassword', 'toggleConfirmPassword');
});
