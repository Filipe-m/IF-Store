package contract

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAddItemToOder(t *testing.T) {

	expectedId := uuid.New().String()
	userId := uuid.New().String()

	go func() {
		app := fiber.New()

		app.Get("/product/:id", func(c *fiber.Ctx) error {
			productId := c.Params("id")
			assert.Equal(t, expectedId, productId)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"product_id": productId, "quantity": 10, "price": 1000})
		})

		err := app.Listen(":9094")
		assert.NoError(t, err)
	}()

	body, _ := json.Marshal(map[string]interface{}{
		"product_id": expectedId,
		"quantity":   30,
	})
	reqOrderItem, err := http.NewRequest(http.MethodPost, "http://localhost:9095/order-item", bytes.NewBuffer(body))
	assert.NoError(t, err)

	reqOrderItem.Header.Set("Content-Type", "application/json")
	reqOrderItem.Header.Set("USER-ID", userId)

	resOrderItem, err := http.DefaultClient.Do(reqOrderItem)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resOrderItem.StatusCode)
}
