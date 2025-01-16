document.addEventListener('DOMContentLoaded', function () {
    var btn = document.querySelector('.button'),
        loader = document.querySelector('.loader'),
        check = document.querySelector('.check');
    
    btn.addEventListener('click', function () {
      loader.classList.add('active');    
    });
   
    loader.addEventListener('animationend', function() {
      check.classList.add('active'); 
    });
  });
  

  // Функция для автоматического добавления пробелов между цифрами при вводе номера карты
document.getElementById("card-number").addEventListener("input", function(e) {
  let value = e.target.value.replace(/\D/g, ''); // Убираем все нецифровые символы
  value = value.replace(/(\d{4})(?=\d)/g, '$1 '); // Добавляем пробелы каждые 4 цифры
  e.target.value = value.slice(0, 19); // Ограничиваем 16 цифрами и пробелами
});

// Функция для автоматического добавления разделителя в дату (MM/YY)
document.getElementById("expiry-date").addEventListener("input", function(e) {
  let value = e.target.value.replace(/\D/g, ''); // Убираем все нецифровые символы
  if (value.length >= 3) {
      value = value.slice(0, 2) + '/' + value.slice(2, 4); // Добавляем / после двух цифр
  }
  e.target.value = value; // Обновляем значение поля
});

// Функция для автоматического добавления пробелов в поле CVV
document.getElementById("cvv").addEventListener("input", function(e) {
  let value = e.target.value.replace(/\D/g, ''); // Убираем все нецифровые символы
  e.target.value = value.slice(0, 3); // Ограничиваем 3 цифрами
});

// Функция для имитации процесса платежа с валидацией
document.getElementById("pay-button").addEventListener("click", function() {
  const cardNumber = document.getElementById('card-number').value.replace(/\D/g, ''); // Убираем все нецифровые символы
  const expiryDate = document.getElementById('expiry-date').value.replace(/\D/g, ''); // Убираем все нецифровые символы
  const cvv = document.getElementById('cvv').value.replace(/\D/g, ''); // Убираем все нецифровые символы
  const cardName = document.getElementById('card-name').value;

  let isValid = true;
  let errorMessage = "";

  // Проверка на номер карты
  if (cardNumber.length !== 16) {
      isValid = false;
      errorMessage += "Card number must have 16 digits.\n";
  }

  // Проверка на дату действия карты (MM/YY)
  if (expiryDate.length !== 4) {
      isValid = false;
      errorMessage += "Expiration date must be in MMYY format.\n";
  }

  // Проверка на CVV (3 цифры)
  if (cvv.length !== 3) {
      isValid = false;
      errorMessage += "CVV must be 3 digits.\n";
  }

  // Проверка на имя на карте
  if (cardName.trim() === "") {
      isValid = false;
      errorMessage += "Name on card cannot be empty.\n";
  }

  // Если есть ошибки, показываем сообщение, иначе начинаем анимацию загрузки
  if (!isValid) {
      document.getElementById('status-message').textContent = errorMessage;
      document.getElementById('payment-status').style.display = 'block';
      document.getElementById('payment-status').style.color = 'red';
  } else {
      // Показываем анимацию на кнопке
      const payButton = document.getElementById('pay-button');
      payButton.disabled = true;
      payButton.querySelector('.spinner-border').style.display = 'inline-block'; // Показываем спиннер
      payButton.querySelector('.fas').style.display = 'none'; // Скрываем стрелку

      // Имитация платежа (задержка)
      setTimeout(function() {
          // После задержки показываем успешный результат и перенаправляем
          document.getElementById('status-message').textContent = "Payment Successful!";
          document.getElementById('payment-status').style.display = 'block';
          document.getElementById('payment-status').style.color = 'green';

          // Перенаправляем на страницу заказов
          window.location.href = '/create-order';
      }, 4000); // Задержка в 3 секунды
  }
});
