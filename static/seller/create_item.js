document.addEventListener("DOMContentLoaded", function () {
    const fileInputs = document.getElementsByClassName("file-upload"); // Множественное число
    const uploadButton = document.getElementById("upload-button");
    const placeholderBox = document.getElementById("placeholder-box");
    const addThumbnailLabel = document.querySelector(".add-thumbnail");
    const thumbnailsContainer = document.getElementById("thumbnails");
    const imageUploadInput = document.querySelector(".add-thumbnail ~ .file-upload"); // Для ввода миниатюр
    const form = document.querySelector("form");

    let uploadedImages = [];

    console.log("DOM loaded, script initialized.");

    // Нажатие на кнопку "Выбрать" для главного изображения
    uploadButton.addEventListener("click", () => {
        console.log("Upload button clicked.");
        fileInputs[0].click();        
    });

    // Обработка главного изображения
    fileInputs[0].addEventListener("change", (event) => {
        console.log("Main image file input change event triggered.");
        const file = event.target.files[0];
        if (file) {
            console.log("Main image file selected:", file.name);
            handleMainFile(file);
            fileInputs[0].value = ""; // Сброс input после выбора
            console.log("File input value reset.");
        } else {
            console.log("No file selected for main image.");
        }
    });

    // Drag and Drop для главного изображения
    placeholderBox.addEventListener("dragover", (event) => {
        console.log("Main image dragover event triggered.");
        event.preventDefault();
        placeholderBox.classList.add("drag-over");
    });

    placeholderBox.addEventListener("dragleave", () => {
        console.log("Main image dragleave event triggered.");
        placeholderBox.classList.remove("drag-over");
    });

    placeholderBox.addEventListener("drop", (event) => {
        console.log("Main image drop event triggered.");
        event.preventDefault();
        placeholderBox.classList.remove("drag-over");
        const file = event.dataTransfer.files[0];
        if (file) {
            console.log("Main image file dropped:", file.name);
            handleMainFile(file);
        } else {
            console.log("No file dropped.");
        }
    });

    // Обработка дополнительных изображений через "+"
    addThumbnailLabel.addEventListener("click", () => {
        console.log("Thumbnail add button clicked.");
        imageUploadInput.click();
    });

    imageUploadInput.addEventListener("change", (event) => {
        console.log("Thumbnail file input change event triggered.");
        const files = Array.from(event.target.files);
        handleThumbnailFiles(files);
        imageUploadInput.value = ""; // Сброс input после выбора
        console.log("Thumbnail file input value reset.");
    });

    // Функция для обработки главного изображения
    function handleMainFile(file) {
        console.log("handleMainFile function triggered.");
        if (file && file.type.startsWith("image/")) {
            const reader = new FileReader();

            reader.onload = function (e) {
                console.log("Main image file successfully read.");
                placeholderBox.innerHTML = `<img src="${e.target.result}" alt="Uploaded Image">`;
                addThumbnailLabel.style.display = "flex"; // Показываем кнопку "+" после загрузки главного изображения
                console.log("Main image displayed, add thumbnail button shown.");
            };

            reader.readAsDataURL(file);
            uploadedImages.push(file);
            console.log("Main image added to uploadedImages array.");
        } else {
            console.log("Invalid file selected for main image.");
            alert("Пожалуйста, выберите изображение.");
        }
    }

    // Функция для обработки миниатюр
    function handleThumbnailFiles(files) {
        console.log("handleThumbnailFiles function triggered.");
        files.forEach((file) => {
            if (uploadedImages.length >= 10) {
                console.log("Maximum of 10 images already uploaded.");
                alert("Можно загрузить не более 10 изображений.");
                return;
            }

            if (file && file.type.startsWith("image/")) {
                console.log("Thumbnail file selected:", file.name);
                const reader = new FileReader();

                reader.onload = function (e) {
                    console.log("Thumbnail file successfully read.");
                    const thumbnail = document.createElement("div");
                    thumbnail.classList.add("thumbnail");

                    thumbnail.innerHTML = `
                        <img src="${e.target.result}" alt="Thumbnail">
                        <button class="delete-thumbnail">&times;</button>
                    `;

                    thumbnail.querySelector(".delete-thumbnail").addEventListener("click", () => {
                        console.log("Delete thumbnail button clicked.");
                        thumbnailsContainer.removeChild(thumbnail);
                        uploadedImages = uploadedImages.filter((img) => img !== file);

                        if (uploadedImages.length < 10) {
                            addThumbnailLabel.style.display = "flex";
                            console.log("Thumbnail deleted, add button shown.");
                        }
                    });

                    thumbnailsContainer.appendChild(thumbnail);
                    console.log("Thumbnail added to thumbnails container.");
                };

                reader.readAsDataURL(file);
                uploadedImages.push(file);
                console.log("Thumbnail added to uploadedImages array.");

                if (uploadedImages.length >= 10) {
                    addThumbnailLabel.style.display = "none";
                    console.log("Maximum number of thumbnails reached, hiding add button.");
                }
            } else {
                console.log("Invalid file selected for thumbnail.");
                alert("Пожалуйста, выберите изображение.");
            }
        });
    }

    // Обработчик отправки формы
    form.addEventListener("submit", function(event) {
        event.preventDefault(); // Предотвращаем стандартное поведение формы

        const formData = new FormData(form); // Получаем данные формы

        // Добавляем изображения в FormData
        uploadedImages.forEach((image) => {
            formData.append("images", image); // Все файлы добавляются под ключом "images"
        });

        // Логируем данные перед отправкой
        const dataToSend = {
            name: formData.get("name"),
            description: formData.get("description"),
            price: parseFloat(formData.get("price")),
            is_available: formData.get("is_available") === "true",
            category: {
                name: formData.get("category")
            }
        };

        // Проверяем значения перед отправкой
        for (const [key, value] of Object.entries(dataToSend)) {
            if (value === undefined || value === null || value === '') {
                console.error(`Ошибка: поле "${key}" не может быть пустым`);
            }
        }

        console.log("Данные перед отправкой на сервер:", dataToSend);

        // Отправляем данные на сервер
        fetch("/seller/products", {
            method: "POST",
            body: formData // Отправляем данные как FormData
        })
        .then(response => {
            if (response.ok) {
                console.log("Product successfully created!");
                return response.json(); // Если сервер возвращает данные
            } else {
                throw new Error(`Unexpected status code: ${response.status}`);
            }
        })
        .then(data => {
            const messageDiv = document.getElementById("message");
            messageDiv.innerText = "Product created successfully!";
            messageDiv.style.color = "green";
            messageDiv.style.display = "block";

            // Очищаем форму
            form.reset();

            // Убираем сообщение через несколько секунд, если нужно
            setTimeout(() => {
                messageDiv.style.display = "none";
            }, 3000);
            window.location.href = "/seller/products"; // Перенаправляем пользователя на страницу с продуктами
        })
        .catch(error => {
            console.error("There was a problem with the fetch operation:", error);
            // Отображаем сообщение об ошибке
            const messageDiv = document.getElementById("message");
            messageDiv.innerText = "There was an error creating the product: " + error.message;
            messageDiv.style.color = "red";
            messageDiv.style.display = "block";
        });
    });
});
