let map; // Глобальная переменная для карты
let clusterer; // Глобальная переменная для кластерера
let selectedPoint = null;
let chooseButton;

function createChooseButton() {
    if (!chooseButton) {
        chooseButton = document.createElement("button");
        chooseButton.textContent = "Choose";
        chooseButton.style.position = "absolute";
        chooseButton.style.zIndex = "1000";
        chooseButton.style.padding = "10px 20px";
        chooseButton.style.backgroundColor = "#007BFF";
        chooseButton.style.color = "white";
        chooseButton.style.border = "none";
        chooseButton.style.borderRadius = "5px";
        chooseButton.style.cursor = "pointer";
        chooseButton.style.display = "none"; // Скрываем по умолчанию
        document.body.appendChild(chooseButton);

        // Добавляем обработчик нажатия кнопки
        chooseButton.addEventListener("click", () => {
            if (selectedPoint) {
                saveSelectedPointToDatabase(selectedPoint);
                alert(`Вы выбрали: ${selectedPoint.name}, ${selectedPoint.address}`);
                chooseButton.style.display = "none"; // Скрываем кнопку после выбора
            }
        });
    }
}

// Показываем кнопку прямо над выбранным маркером
function showChooseButton(coords) {
    const projection = map.options.get("projection"); // Используем проекцию карты
    const position = map.converter.globalToPage(
        projection.toGlobalPixels(coords, map.getZoom()) // Преобразуем координаты в пиксели на экране
    );

    // Устанавливаем положение кнопки в пиксельных координатах
    chooseButton.style.left = `${position[0] - chooseButton.offsetWidth / 2}px`;
    chooseButton.style.top = `${position[1] - chooseButton.offsetHeight - 10}px`;
    chooseButton.style.display = "block"; // Показываем кнопку
}

async function findNearestPoint() {
    const addressInput = document.getElementById("search-address").value;

    if (!addressInput) {
        alert("Введите адрес для поиска!");
        return;
    }

    try {
        const geocodeResponse = await ymaps.geocode(addressInput);
        const geoObject = geocodeResponse.geoObjects.get(0);

        if (!geoObject) {
            alert("Адрес не найден. Проверьте правильность ввода.");
            return;
        }

        const addressCoords = geoObject.geometry.getCoordinates();

        const response = await fetch("/api/pickup-points");
        const points = await response.json();

        if (!points || points.length === 0) {
            alert("Нет доступных пунктов выдачи.");
            return;
        }

        let nearestPoint = null;
        let minDistance = Infinity;

        points.forEach((point) => {
            const distance = ymaps.coordSystem.geo.getDistance(
                addressCoords,
                [point.latitude, point.longitude]
            );

            if (distance < minDistance) {
                minDistance = distance;
                nearestPoint = point;
            }
        });

        if (nearestPoint) {
            map.setCenter([nearestPoint.latitude, nearestPoint.longitude], 15);


            const placemark = new ymaps.Placemark(
                [nearestPoint.latitude, nearestPoint.longitude],
                {
                    balloonContent: `<b>${nearestPoint.name}</b><br>${nearestPoint.address}`,
                },
                {
                    iconLayout: "default#image",
                    iconImageHref: "/assets/logo.png",
                    iconImageSize: [40, 40],
                    iconImageOffset: [-20, -20],
                }
            );

            placemark.events.add("click", () => {
                selectedPoint = nearestPoint; // Устанавливаем выбранную точку
                showChooseButton([nearestPoint.latitude, nearestPoint.longitude]);
            });

            clusterer.add(placemark);
        }
    } catch (error) {
        console.error("Ошибка при поиске ближайшего пункта:", error);
        alert("Произошла ошибка при поиске. Попробуйте снова.");
    }
}

// Функция для очистки поиска
function clearSearch() {
    const searchInput = document.getElementById("search-address");
    searchInput.value = ""; // Очистить поле ввода
    document.getElementById("clear-search").style.display = "none"; // Скрыть крестик

    // Загрузка всех пунктов выдачи обратно
    fetch("/api/pickup-points")
        .then((response) => response.json())
        .then((points) => {
            // Сбрасываем карту и боковую панель
            clusterer.removeAll();
            const geoObjects = points.map((point) => {
                return new ymaps.Placemark(
                    [point.latitude, point.longitude],
                    {
                        balloonContent: `<b>${point.name}</b><br>${point.address}`,
                    },
                    {
                        iconLayout: "default#image",
                        iconImageHref: "/assets/logo.png",
                        iconImageSize: [40, 40],
                        iconImageOffset: [-20, -20],
                    }
                );
            });

            clusterer.add(geoObjects);

            // Обновляем список всех пунктов выдачи на боковой панели
            updateSidebar(points);
        })
        .catch((error) =>
            console.error("Ошибка загрузки пунктов выдачи:", error)
        );
}

// Инициализация карты
document.addEventListener("DOMContentLoaded", function () {
    ymaps.ready(initMap);
    createChooseButton();
    // Добавляем обработчик Enter для поиска
    const searchInput = document.getElementById("search-address");
    searchInput.addEventListener("keypress", function (event) {
        if (event.key === "Enter") {
            event.preventDefault();
            findNearestPoint(); // Запуск поиска
        }
    });

    // Добавляем обработчик для крестика
    document.getElementById("clear-search").addEventListener("click", clearSearch);
});

function initMap() {
    if (!map) {
        map = new ymaps.Map("map", {
            center: [43.25667, 76.92861],
            zoom: 12,
            controls: ["zoomControl"],
        });

        clusterer = new ymaps.Clusterer();
        map.geoObjects.add(clusterer);

        // Загрузка всех пунктов выдачи
        fetch("/api/pickup-points")
            .then((response) => response.json())
            .then((points) => {
                const geoObjects = points.map((point) => {
                    const placemark = new ymaps.Placemark(
                        [point.latitude, point.longitude],
                        {
                            balloonContent: `<b>${point.name}</b><br>${point.address}`,
                        },
                        {
                            iconLayout: "default#image",
                            iconImageHref: "/assets/logo.png",
                            iconImageSize: [40, 40],
                            iconImageOffset: [-20, -20],
                        }
                    );

                    // Показать кнопку при нажатии на метку
                    placemark.events.add("click", () => {
                        selectedPoint = point;
                        showChooseButton([point.latitude, point.longitude]);
                    });

                    return placemark;
                });

                clusterer.add(geoObjects);

                // Обновляем список пунктов на боковой панели
                updateSidebar(points);
            })
            .catch((error) => console.error("Ошибка загрузки пунктов выдачи:", error));
    }
}

function saveSelectedPointToDatabase(point) {
    console.log("Attempting to save selected point to database...");
    console.log("Point data:", point);

    fetch('/api/save-delivery-address', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            user_id: 1, // Замените на текущий ID пользователя (например, из сессии)
            address_line: point.address,
            city: "Алматы", // Можно динамически получить, если нужно
            latitude: point.latitude,
            longitude: point.longitude
        }),
    })
        .then((response) => {
            console.log("Response received from server:", response);
            if (!response.ok) {
                console.error("Server responded with an error:", response.status, response.statusText);
                throw new Error('Ошибка при сохранении адреса.');
            }
            return response.json();
        })
        .then((data) => {
            console.log("Server responded with data:", data);
            alert('Адрес успешно сохранён!');
        })
        .catch((error) => {
            console.error("Ошибка при сохранении адреса:", error);
            alert('Не удалось сохранить адрес.');
        });
}


// Добавляем сохранение при клике на пункт в боковой панели
function updateSidebar(points) {
    const sidebar = document.getElementById("sidebar");
    sidebar.innerHTML = "";

    points.forEach((point) => {
        const item = document.createElement("div");
        item.className = "sidebar-item";
        item.style.padding = "10px";
        item.style.borderBottom = "1px solid #ddd";
        item.style.cursor = "pointer";

        item.innerHTML = `<b>${point.name}</b><br>${point.address}`;

        // Центрирование карты и показ кнопки при клике
        item.addEventListener("click", () => {
            map.setCenter([point.latitude, point.longitude], 15);
            selectedPoint = point;
            showChooseButton([point.latitude, point.longitude]);
        });

        sidebar.appendChild(item);
    });
}


function goBack() {
    window.history.back(); // Возвращает пользователя на предыдущую страницу
}

