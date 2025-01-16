package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/paymentintent"
	"net/http"
)

type CreatePaymentIntentResponse struct {
	ClientSecret string `json:"clientSecret"`
}

func CreatePaymentIntentHandler(c *gin.Context) {
	// Если это GET запрос, то возвращаем checkout.html
	if c.Request.Method == http.MethodGet {
		// Отправляем файл checkout.html
		c.HTML(http.StatusOK, "checkout.html", nil)
		return
	}

	// Если это POST запрос, создаем PaymentIntent
	if c.Request.Method == http.MethodPost {
		// Устанавливаем Stripe секретный ключ
		stripe.Key = "sk_test_51QQuKLGWZ73ZM8VUckKCyVUcTqoSXHzSOJeMtb9ioaOfNLdjGZtByMOuTMZFXB5CGDgFcZNOOe6atfd98T0LgnTX00bt1BI9N8" // Вставьте свой секретный ключ

		// Создаем PaymentIntent
		params := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(1000), // Сумма в центах (например, $10.00)
			Currency: stripe.String("usd"),
		}

		intent, err := paymentintent.New(params)
		if err != nil {
			// Обрабатываем ошибку и отправляем 500
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Формируем ответ с client_secret
		response := CreatePaymentIntentResponse{
			ClientSecret: intent.ClientSecret,
		}

		// Отправляем ответ в формате JSON
		c.JSON(http.StatusOK, response)
		return
	}

	// Если это другой метод (например, PUT или DELETE), отправляем ошибку
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"error": "Method not allowed",
	})
}
