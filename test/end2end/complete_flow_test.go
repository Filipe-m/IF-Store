package end2end

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"net/http"
	"testing"
)

func TestCompleteFlow(t *testing.T) {
	randomString := uuid.New().String()

	username := "example_user_" + randomString
	email := "user_" + randomString + "@example.com"
	productName := "example_product_" + randomString
	productDescription := "example product description " + randomString
	price := int(math.Floor(rand.Float64()*1000) + 1)
	quantity := int(math.Floor(rand.Float64()*100) + 1)

	bodyUserJSON, _ := json.Marshal(map[string]string{
		"username": username,
		"email":    email,
	})

	resUser, err := http.Post("http://localhost:9091/users", "application/json", bytes.NewBuffer(bodyUserJSON))
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resUser.StatusCode)

	var data map[string]interface{}
	json.NewDecoder(resUser.Body).Decode(&data)
	userID := data["id"].(string)

	resPaymentMethods, err := http.Get("http://localhost:9096/paymentMethod/" + userID)
	assert.NoError(t, err)

	var dataMethods []map[string]interface{}
	err = json.NewDecoder(resPaymentMethods.Body).Decode(&dataMethods)
	assert.NoError(t, err)
	paymentMethodId := dataMethods[0]["id"].(string)

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
	reqOrderItem.Header.Set("USER-ID", userID)

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
	reqFinishOrder.Header.Set("USER-ID", userID)

	resFinishOrder, err := http.DefaultClient.Do(reqFinishOrder)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resFinishOrder.StatusCode)
}

func BenchmarkCompleteFlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomString := uuid.New().String()

		username := "example_user_" + randomString
		email := "user_" + randomString + "@example.com"
		productName := "example_product_" + randomString
		productDescription := "example product description " + randomString
		price := int(math.Floor(rand.Float64()*1000) + 1)
		quantity := int(math.Floor(rand.Float64()*100) + 1)

		bodyUserJSON, _ := json.Marshal(map[string]string{
			"username": username,
			"email":    email,
		})

		resUser, err := http.Post("http://localhost:9091/users", "application/json", bytes.NewBuffer(bodyUserJSON))
		assert.NoError(b, err)

		assert.Equal(b, http.StatusCreated, resUser.StatusCode)

		var data map[string]interface{}
		json.NewDecoder(resUser.Body).Decode(&data)
		userID := data["id"].(string)

		bodyProductJSON, _ := json.Marshal(map[string]interface{}{
			"name":        productName,
			"description": productDescription,
			"price":       price,
		})
		resProduct, err := http.Post("http://localhost:9094/product/register", "application/json", bytes.NewBuffer(bodyProductJSON))
		assert.NoError(b, err)

		assert.Equal(b, http.StatusCreated, resProduct.StatusCode)

		json.NewDecoder(resProduct.Body).Decode(&data)
		productID := data["id"].(string)

		bodyStockJSON, _ := json.Marshal(map[string]interface{}{
			"quantity": quantity,
		})
		reqStock, err := http.NewRequest(http.MethodPut, "http://localhost:9094/stock/"+productID+"/add", bytes.NewBuffer(bodyStockJSON))
		assert.NoError(b, err)
		reqStock.Header.Set("Content-Type", "application/json")

		resStock, err := http.DefaultClient.Do(reqStock)
		assert.NoError(b, err)

		assert.Equal(b, http.StatusCreated, resStock.StatusCode)

		bodyOrderItemJSON, _ := json.Marshal(map[string]interface{}{
			"product_id": productID,
			"quantity":   quantity,
		})
		reqOrderItem, err := http.NewRequest(http.MethodPost, "http://localhost:9095/order-item", bytes.NewBuffer(bodyOrderItemJSON))
		assert.NoError(b, err)

		reqOrderItem.Header.Set("Content-Type", "application/json")
		reqOrderItem.Header.Set("USER-ID", userID)

		resOrderItem, err := http.DefaultClient.Do(reqOrderItem)
		assert.NoError(b, err)
		assert.Equal(b, http.StatusCreated, resOrderItem.StatusCode)

		json.NewDecoder(resOrderItem.Body).Decode(&data)
		orderID := data["id"].(string)

		bodyFinishOrderJSON, _ := json.Marshal(map[string]interface{}{
			"order_id":     orderID,
			"payment_data": "lorem",
		})
		reqFinishOrder, err := http.NewRequest(http.MethodPost, "http://localhost:9095/order/finish", bytes.NewBuffer(bodyFinishOrderJSON))
		assert.NoError(b, err)

		reqFinishOrder.Header.Set("Content-Type", "application/json")
		reqFinishOrder.Header.Set("USER-ID", userID)

		resFinishOrder, err := http.DefaultClient.Do(reqFinishOrder)
		assert.NoError(b, err)
		assert.Equal(b, http.StatusCreated, resFinishOrder.StatusCode)
	}
}
