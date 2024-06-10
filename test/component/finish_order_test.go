package component

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"net/http"
	"testing"
)

func TestCompleteFlow(t *testing.T) {

	go func() {
		app := fiber.New()

		var executionCount int

		app.Post("/send-message", func(c *fiber.Ctx) error {
			executionCount++

			var request map[string]interface{}
			err := c.BodyParser(&request)
			assert.NoError(t, err)

			switch executionCount {
			case 1:
				assert.Equal(t, "Your order has been placed", request["message"])
			case 2:
				assert.Equal(t, "Your order has been shipped", request["message"])
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		})

		err := app.Listen(":9092")
		assert.NoError(t, err)
	}()

	randomString := uuid.New().String()

	productName := "example_product_" + randomString
	productDescription := "example product description " + randomString
	price := int(math.Floor(rand.Float64()*1000) + 1)
	quantity := int(math.Floor(rand.Float64()*100) + 1)

	resPaymentMethods, err := http.Get("http://localhost:9096/paymentMethod/" + randomString)
	assert.NoError(t, err)

	var dataMethods []map[string]interface{}
	err = json.NewDecoder(resPaymentMethods.Body).Decode(&dataMethods)
	assert.NoError(t, err)
	paymentMethodId := dataMethods[0]["id"].(string)

	var data map[string]interface{}

	bodyProductJSON, _ := json.Marshal(map[string]interface{}{
		"name":        productName,
		"description": productDescription,
		"price":       price,
	})
	resProduct, err := http.Post("http://localhost:9094/product/register", "application/json", bytes.NewBuffer(bodyProductJSON))
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resProduct.StatusCode)

	json.NewDecoder(resProduct.Body).Decode(&data)
	productID := data["id"].(string)

	bodyStockJSON, _ := json.Marshal(map[string]interface{}{
		"quantity": quantity,
	})
	reqStock, err := http.NewRequest(http.MethodPut, "http://localhost:9094/stock/"+productID+"/add", bytes.NewBuffer(bodyStockJSON))
	assert.NoError(t, err)
	reqStock.Header.Set("Content-Type", "application/json")

	resStock, err := http.DefaultClient.Do(reqStock)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resStock.StatusCode)

	bodyOrderItemJSON, _ := json.Marshal(map[string]interface{}{
		"product_id": productID,
		"quantity":   quantity,
	})
	reqOrderItem, err := http.NewRequest(http.MethodPost, "http://localhost:9095/order-item", bytes.NewBuffer(bodyOrderItemJSON))
	assert.NoError(t, err)

	reqOrderItem.Header.Set("Content-Type", "application/json")
	reqOrderItem.Header.Set("USER-ID", randomString)

	resOrderItem, err := http.DefaultClient.Do(reqOrderItem)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resOrderItem.StatusCode)

	json.NewDecoder(resOrderItem.Body).Decode(&data)
	orderID := data["id"].(string)

	bodyFinishOrderJSON, _ := json.Marshal(map[string]interface{}{
		"order_id":          orderID,
		"payment_method_id": paymentMethodId,
	})
	reqFinishOrder, err := http.NewRequest(http.MethodPost, "http://localhost:9095/order/finish", bytes.NewBuffer(bodyFinishOrderJSON))
	assert.NoError(t, err)

	reqFinishOrder.Header.Set("Content-Type", "application/json")
	reqFinishOrder.Header.Set("USER-ID", randomString)

	resFinishOrder, err := http.DefaultClient.Do(reqFinishOrder)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resFinishOrder.StatusCode)
}
