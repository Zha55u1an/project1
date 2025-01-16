// authentication.go

package handlers

import (
	"fmt"
	"go_project/internal/models"
	"go_project/pkg/db"
	"go_project/pkg/utils"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

var jwtKey = []byte("jwtkey_go_project")

func Login(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		// Если метод GET, отображаем страницу login.html
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}

	var user models.User

	// Парсим JSON из POST-запроса
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User

	// Проверяем существование пользователя
	db.DB.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID == 0 {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	// Проверяем пароль
	errHash := utils.CheckPasswordHash(user.Password, existingUser.Password)
	if !errHash {
		c.JSON(400, gin.H{"error": "invalid password"})
		return
	}

	// Генерация токена
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &models.Claims{
		Role: existingUser.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.Email,
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not generate token"})
		return
	}

	// Устанавливаем cookie с токеном
	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)

	// Указываем URL перенаправления в зависимости от роли
	var redirectURL string
	if existingUser.Role == "admin" {
		redirectURL = "/admin"
	} else if existingUser.Role == "buyer" {
		redirectURL = "/home"
	} else if existingUser.Role == "seller" {
		redirectURL = "/seller/dashboard"
	} else {
		c.JSON(400, gin.H{"error": "invalid role"})
		return
	}

	c.JSON(200, gin.H{"success": "user logged in", "redirect": redirectURL})
}



func Signup(c *gin.Context) {
    if c.Request.Method == http.MethodGet {
        log.Println("GET request received for /signup")
        c.HTML(http.StatusOK, "signup.html", nil)
        return
    }

    log.Println("POST request received for /signup")
    var request struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
        Role     string `json:"role" binding:"required"`
    }

    // Шаг 1: Получаем данные из запроса
    if err := c.ShouldBindJSON(&request); err != nil {
        log.Printf("[ERROR] Error binding JSON: %v", err)
        c.JSON(400, gin.H{"error": "Invalid request data. Please check your input."})
        return
    }
    log.Printf("[INFO] Request data: Email=%s, Role=%s", request.Email, request.Role)

    // Шаг 2: Проверяем, существует ли пользователь
    var existingUser models.User
    if err := db.DB.Where("email = ?", request.Email).First(&existingUser).Error; err == nil && existingUser.ID != 0 {
        log.Printf("[ERROR] User with email %s already exists", request.Email)
        c.JSON(400, gin.H{"error": "User with this email already exists"})
        return
    }

    // Шаг 3: Хешируем пароль
    hashedPassword, err := utils.PasswordHash(request.Password)
    if err != nil {
        log.Printf("[ERROR] Error hashing password: %v", err)
        c.JSON(500, gin.H{"error": "Failed to hash password"})
        return
    }
    log.Println("[INFO] Password hashed successfully")

    // Шаг 4: Генерируем код подтверждения
    code := GenerateVerificationCode()
    log.Printf("[INFO] Generated verification code: %s for email %s", code, request.Email)

    // Шаг 5: Сохраняем или обновляем запись в Verification
    var verification models.Verification
    if err := db.DB.Where("email = ?", request.Email).First(&verification).Error; err == nil {
        // Если запись существует, обновляем код, пароль и время истечения
        verification.VerificationCode = code
        verification.Password = hashedPassword
        verification.Role = request.Role
        verification.ExpiresAt = time.Now().Add(15 * time.Minute)
        if err := db.DB.Save(&verification).Error; err != nil {
            log.Printf("[ERROR] Error updating verification record: %v", err)
            c.JSON(500, gin.H{"error": "Failed to update verification record"})
            return
        }
        log.Printf("[INFO] Verification record updated for email %s", request.Email)
    } else {
        // Если записи нет, создаём новую
        verification = models.Verification{
            Email:            request.Email,
            VerificationCode: code,
            Password:         hashedPassword,
            Role:             request.Role,
            ExpiresAt:        time.Now().Add(15 * time.Minute),
        }
        if err := db.DB.Create(&verification).Error; err != nil {
            log.Printf("[ERROR] Error saving verification record: %v", err)
            c.JSON(500, gin.H{"error": "Failed to save verification record"})
            return
        }
        log.Printf("[INFO] Verification record created for email %s", request.Email)
    }

    // Шаг 6: Отправляем код подтверждения на email
    if err := SendVerificationEmail(request.Email, code); err != nil {
        log.Printf("[ERROR] Error sending verification email: %v", err)
        c.JSON(500, gin.H{"error": "Failed to send verification email"})
        return
    }
    log.Printf("[INFO] Verification email sent to %s", request.Email)

    // Шаг 7: Возвращаем успешный ответ для клиента
    c.JSON(200, gin.H{
        "message": "Verification email sent. Please check your inbox.",
        "email":   request.Email,
    })
}








func Logout(c *gin.Context) {
	// Удаляем cookie с токеном
	c.SetCookie("token", "", -1, "/", "localhost", false, true)

	// Логируем успешный выход
	log.Println("User logged out successfully")

	// Перенаправляем пользователя на страницу входа
	c.Redirect(http.StatusFound, "/login")
}


func GenerateVerificationCode() string {
    log.Println("Generating verification code")
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    code := fmt.Sprintf("%06d", rng.Intn(1000000))
    log.Printf("Generated verification code: %s", code)
    return code
}


func SendVerificationEmail(to string, code string) error {
    log.Printf("Preparing to send email to %s", to)
    mailer := gomail.NewMessage()
    mailer.SetHeader("From", "your-email@gmail.com")
    mailer.SetHeader("To", to)
    mailer.SetHeader("Subject", "Email Verification Code")
    mailer.SetBody("text/plain", "Your verification code is: "+code)

    dialer := gomail.NewDialer("smtp.gmail.com", 587, "zhasikshokan90@gmail.com", "oarj dmkh nnsr ttmx")

    log.Printf("Sending email to %s with code %s", to, code)
    err := dialer.DialAndSend(mailer)
    if err != nil {
        log.Printf("Failed to send email: %v", err)
        return err
    }
    log.Printf("Email successfully sent to %s", to)
    return nil
}

func ResendCode(c *gin.Context) {
    log.Println("POST request received for /resend-code")
    var request struct {
        Email string `json:"email"`
    }

    // Получаем email из запроса
    if err := c.ShouldBindJSON(&request); err != nil {
        log.Printf("Error binding JSON: %v", err)
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    log.Printf("Resend request received for email: %s", request.Email)

    // Проверяем, существует ли email
    var verification models.Verification
    if err := db.DB.Where("email = ?", request.Email).First(&verification).Error; err != nil {
        log.Printf("No verification process found for email %s: %v", request.Email, err)
        c.JSON(400, gin.H{"error": "No verification process found for this email"})
        return
    }

    // Генерируем новый код
    code := GenerateVerificationCode()
    log.Printf("Generated new verification code: %s for email %s", code, request.Email)

    // Обновляем код в базе данных
    verification.VerificationCode = code
    verification.ExpiresAt = time.Now().Add(15 * time.Minute)
    if err := db.DB.Save(&verification).Error; err != nil {
        log.Printf("Error updating verification code in database: %v", err)
        c.JSON(500, gin.H{"error": "Failed to update verification code"})
        return
    }
    log.Printf("Verification code updated in database for email %s", request.Email)

    // Отправляем код
    if err := SendVerificationEmail(request.Email, code); err != nil {
        log.Printf("Failed to resend verification email to %s: %v", request.Email, err)
        c.JSON(500, gin.H{"error": "Failed to resend verification email"})
        return
    }
    log.Printf("Verification email resent to %s", request.Email)

    c.JSON(200, gin.H{"message": "Verification email resent successfully"})
}

func Verification(c *gin.Context) {
    if c.Request.Method == http.MethodGet {
        log.Println("[INFO] GET request received for /verification")
        
        // Получаем email из параметра запроса
        email := c.Query("email")
        if email == "" {
            log.Println("[ERROR] Missing email in GET /verification")
            c.HTML(http.StatusBadRequest, "error.html", gin.H{
                "message": "Email is required to access verification page.",
            })
            return
        }
        
        log.Printf("[INFO] Rendering verification page for email: %s", email)
        c.HTML(http.StatusOK, "verification.html", gin.H{
            "email": email,
        })
        return
    }

    log.Println("[INFO] POST request received for /verification")

    // Структура для получения данных из JSON-запроса
    var request struct {
        Email string `json:"email" binding:"required,email"`
        Code  string `json:"code" binding:"required"`
    }

    // Привязка JSON данных к структуре
    if err := c.ShouldBindJSON(&request); err != nil {
        log.Printf("[ERROR] Error binding JSON: %v", err)
        c.JSON(400, gin.H{"error": "Invalid request data. Please check your input."})
        return
    }
    log.Printf("[INFO] Parsed request: Email=%s, Code=%s", request.Email, request.Code)

    // Шаг 1: Проверяем наличие записи с email и кодом
    var verification models.Verification
    if err := db.DB.Where("email = ? AND verification_code = ?", request.Email, request.Code).First(&verification).Error; err != nil {
        log.Printf("[ERROR] Invalid or expired verification code for email %s", request.Email)
        c.JSON(400, gin.H{"error": "Invalid or expired verification code"})
        return
    }
    log.Printf("[INFO] Verification record found: Email=%s, Code=%s, ExpiresAt=%v", verification.Email, verification.VerificationCode, verification.ExpiresAt)

    // Шаг 2: Проверяем срок действия кода
    if time.Now().After(verification.ExpiresAt) {
        log.Printf("[ERROR] Verification code expired for email %s", request.Email)
        db.DB.Delete(&verification) // Удаляем устаревший код
        log.Printf("[INFO] Expired verification record deleted for email %s", request.Email)
        c.JSON(400, gin.H{"error": "Verification code has expired. Please request a new one."})
        return
    }

    // Шаг 3: Создаём пользователя
    newUser := models.User{
        Email:      verification.Email,
        Password:   verification.Password, // Пароль сохранён на этапе регистрации
        Role:       verification.Role,    // Используем роль из Verification
        IsVerified: true,
    }
    log.Printf("[INFO] Creating new user: Email=%s, Role=%s", newUser.Email, newUser.Role)

    if err := db.DB.Create(&newUser).Error; err != nil {
        log.Printf("[ERROR] Error creating user for email %s: %v", request.Email, err)
        c.JSON(500, gin.H{"error": "Failed to create user. Please try again later."})
        return
    }
    log.Printf("[INFO] User created successfully: Email=%s, Role=%s", newUser.Email, newUser.Role)

    // Шаг 4: Удаляем запись из таблицы Verification
    if err := db.DB.Unscoped().Delete(&verification).Error; err != nil {
        log.Printf("[ERROR] Error deleting verification record for email %s: %v", request.Email, err)
        c.JSON(500, gin.H{"error": "Failed to clean up verification record."})
        return
    }
    log.Printf("[INFO] Verification record deleted for email %s", request.Email)

    // Шаг 5: Возвращаем успешный ответ
    log.Printf("[INFO] Verification successful for email %s. Redirecting to /login", request.Email)
    c.JSON(200, gin.H{
        "message":  "Verification successful. Please log in.",
        "redirect": "/login",
    })
}

func VerifyResetCode(c *gin.Context) {
    if c.Request.Method == http.MethodGet {
        log.Println("[INFO] GET request received for /verify-reset-code")

        // Получаем email из параметра запроса
        email := c.Query("email")
        if email == "" {
            log.Println("[ERROR] Missing email in GET /verify-reset-code")
            c.HTML(http.StatusBadRequest, "error.html", gin.H{
                "message": "Email is required to access reset code verification page.",
            })
            return
        }

        log.Printf("[INFO] Rendering reset code verification page for email: %s", email)
        c.HTML(http.StatusOK, "verify-reset-code.html", gin.H{
            "email": email,
        })
        return
    }

    log.Println("[INFO] POST request received for /verify-reset-code")

    // Структура для получения данных из JSON-запроса
    var request struct {
        Email string `json:"email" binding:"required,email"`
        Code  string `json:"code" binding:"required"`
    }

    // Привязка JSON данных к структуре
    if err := c.ShouldBindJSON(&request); err != nil {
        log.Printf("[ERROR] Error binding JSON: %v", err)
        c.JSON(400, gin.H{"error": "Invalid request data. Please check your input."})
        return
    }
    log.Printf("[INFO] Parsed request: Email=%s, Code=%s", request.Email, request.Code)

    // Шаг 1: Проверяем наличие записи с email и кодом
    var verification models.Verification
    if err := db.DB.Where("email = ? AND verification_code = ?", request.Email, request.Code).First(&verification).Error; err != nil {
        log.Printf("[ERROR] Invalid or expired reset code for email %s", request.Email)
        c.JSON(400, gin.H{"error": "Invalid or expired reset code"})
        return
    }
    log.Printf("[INFO] Verification record found: Email=%s, Code=%s, ExpiresAt=%v", verification.Email, verification.VerificationCode, verification.ExpiresAt)

    // Шаг 2: Проверяем срок действия кода
    if time.Now().After(verification.ExpiresAt) {
        log.Printf("[ERROR] Reset code expired for email %s", request.Email)
        db.DB.Delete(&verification) // Удаляем устаревший код
        log.Printf("[INFO] Expired reset code record deleted for email %s", request.Email)
        c.JSON(400, gin.H{"error": "Reset code has expired. Please request a new one."})
        return
    }

    // Шаг 3: Перенаправляем пользователя на страницу обновления пароля
    log.Printf("[INFO] Reset code verified for email %s. Redirecting to /update-password", request.Email)
    c.JSON(200, gin.H{
        "message":  "Reset code verified. Redirecting to update password page.",
        "redirect": "/update-password?email=" + request.Email,
    })
}


func ResetPassword(c *gin.Context) {
    if c.Request.Method == http.MethodGet {
        log.Println("[INFO] GET request received for /reset-password")
        c.HTML(http.StatusOK, "reset-password.html", nil)
        return
    }

    log.Println("[INFO] POST request received for /reset-password")

    // Структура для получения данных из JSON-запроса
    var request struct {
        Email string `json:"email" binding:"required,email"`
    }

    // Привязка JSON данных
    if err := c.ShouldBindJSON(&request); err != nil {
        log.Printf("[ERROR] Error binding JSON: %v", err)
        c.JSON(400, gin.H{"error": "Invalid email format."})
        return
    }

    log.Printf("[INFO] Reset password request for email: %s", request.Email)

    // Проверяем, существует ли пользователь с этим email
    var user models.User
    if err := db.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
        log.Printf("[ERROR] Email not found: %s", request.Email)
        c.JSON(404, gin.H{"error": "Email not registered."})
        return
    }

    // Генерируем код восстановления
    resetCode := GenerateVerificationCode()
    log.Printf("[INFO] Generated reset code for email %s: %s", request.Email, resetCode)

    // Сохраняем или обновляем код в таблице Verification
    var verification models.Verification
    if err := db.DB.Where("email = ?", request.Email).First(&verification).Error; err == nil {
        // Если запись уже существует, обновляем код восстановления и время действия
        log.Printf("[INFO] Existing verification record found for email: %s", request.Email)
        if err := db.DB.Model(&verification).Updates(models.Verification{
            VerificationCode: resetCode,
            ExpiresAt:        time.Now().Add(15 * time.Minute),
        }).Error; err != nil {
            log.Printf("[ERROR] Error updating reset code: %v", err)
            c.JSON(500, gin.H{"error": "Failed to generate reset code. Please try again."})
            return
        }
        log.Printf("[INFO] Reset code updated for email: %s", request.Email)
    } else if err == gorm.ErrRecordNotFound {
        // Если записи нет, создаём новую
        log.Printf("[INFO] No existing verification record found. Creating new for email: %s", request.Email)
        verification = models.Verification{
            Email:            request.Email,
            VerificationCode: resetCode,
            ExpiresAt:        time.Now().Add(15 * time.Minute),
        }
        if err := db.DB.Create(&verification).Error; err != nil {
            log.Printf("[ERROR] Error saving reset code: %v", err)
            c.JSON(500, gin.H{"error": "Failed to generate reset code. Please try again."})
            return
        }
        log.Printf("[INFO] Reset code created for email: %s", request.Email)
    } else {
        // Если произошла другая ошибка, логируем её
        log.Printf("[ERROR] Unexpected error querying verification record: %v", err)
        c.JSON(500, gin.H{"error": "Failed to process reset request. Please try again."})
        return
    }

    // Отправляем код на email
    if err := SendVerificationEmail(request.Email, resetCode); err != nil {
        log.Printf("[ERROR] Failed to send reset code to email: %s. Error: %v", request.Email, err)
        c.JSON(500, gin.H{"error": "Failed to send reset code. Please try again."})
        return
    }

    log.Printf("[INFO] Reset code sent to email: %s", request.Email)

    // Возвращаем JSON-ответ с редиректом
    c.JSON(200, gin.H{
        "message":  "Reset code sent successfully.",
        "redirect": "/verify-reset-code?email=" + request.Email,
    })
}




func UpdatePassword(c *gin.Context) {
    if c.Request.Method == http.MethodGet {
        log.Println("[INFO] GET request received for /update-password")
        
        email := c.Query("email")
        if email == "" {
            log.Println("[ERROR] Missing email in GET /update-password")
            c.HTML(http.StatusBadRequest, "error.html", gin.H{
                "message": "Email is required to update the password.",
            })
            return
        }

        log.Printf("[INFO] Rendering password update page for email: %s", email)
        c.HTML(http.StatusOK, "update-password.html", gin.H{
            "email": email,
        })
        return
    }

    log.Println("[INFO] POST request received for /update-password")

    var request struct {
        Email        string `json:"email" binding:"required,email"`
        NewPassword  string `json:"new_password" binding:"required"`
        ConfirmPassword string `json:"confirm_password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        log.Printf("[ERROR] Error binding JSON: %v", err)
        c.JSON(400, gin.H{"error": "Invalid input data. Please check your inputs."})
        return
    }

    // Проверка совпадения паролей
    if request.NewPassword != request.ConfirmPassword {
        log.Printf("[ERROR] Passwords do not match for email: %s", request.Email)
        c.JSON(400, gin.H{"error": "Passwords do not match."})
        return
    }

    log.Printf("[INFO] Updating password for email: %s", request.Email)

    // Хешируем новый пароль
    hashedPassword, err := utils.PasswordHash(request.NewPassword)
    if err != nil {
        log.Printf("[ERROR] Failed to hash password: %v", err)
        c.JSON(500, gin.H{"error": "Failed to update password. Please try again later."})
        return
    }

    // Обновляем пароль в базе данных
    if err := db.DB.Model(&models.User{}).Where("email = ?", request.Email).Update("password", hashedPassword).Error; err != nil {
        log.Printf("[ERROR] Failed to update password for email: %s", request.Email)
        c.JSON(500, gin.H{"error": "Failed to update password. Please try again later."})
        return
    }

    log.Printf("[INFO] Password updated successfully for email: %s", request.Email)
    c.JSON(200, gin.H{"message": "Password updated successfully. You can now log in."})
}


