package handlers

import (
    "fmt"
    "log"
    "strings"
    "github.com/google/uuid"
    "github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
    // Получаем файл из запроса
    file, err := c.FormFile("image")
    if err != nil {
        log.Println("Ошибка при загрузке изображения:", err)
        c.JSON(500, gin.H{"status": 500, "message": "Ошибка сервера"})
        return
    }

    // Генерация уникального имени файла
    uniqueId := uuid.New()
    filename := strings.Replace(uniqueId.String(), "-", "", -1)

    // Определение расширения файла
    fileExt := strings.Split(file.Filename, ".")[1]
    imageName := fmt.Sprintf("%s.%s", filename, fileExt)

    // Путь для сохранения изображения
    filePath := fmt.Sprintf("./uploads/%s", imageName)

    // Сохранение файла на сервере
    if err := c.SaveUploadedFile(file, filePath); err != nil {
        log.Println("Ошибка при сохранении изображения:", err)
        c.JSON(500, gin.H{"status": 500, "message": "Ошибка сервера"})
        return
    }

    // Формирование URL для изображения
    imageUrl := fmt.Sprintf("http://localhost:8080/uploads/%s", imageName)

    // Возвращаем успешный ответ с данными об изображении
    data := map[string]interface{}{
        "imageName": imageName,
        "imageUrl":  imageUrl,
        "header":    file.Header,
        "size":      file.Size,
    }

    c.JSON(201, gin.H{"status": 201, "message": "Изображение успешно загружено", "data": data})
}
