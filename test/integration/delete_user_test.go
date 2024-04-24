package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"test/platform"
	"testing"
	"time"
)

func TestDeleteUserIntegration(t *testing.T) {

	username := fmt.Sprintf("example_user_%s", uuid.New().String())
	email := fmt.Sprintf("user_%s@example.com", uuid.New().String())

	body := map[string]string{
		"username": username,
		"email":    email,
	}

	jsonBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:9091/users", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)

	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)

	userId, ok := response["id"].(string)
	assert.True(t, ok)

	req, err = http.NewRequest(http.MethodDelete, "http://localhost:9091/users/"+userId, bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	res, err = client.Do(req)
	assert.NoError(t, err)

	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	connect, err := platform.NewPostgresConnect(platform.Database{
		Localhost: "localhost",
		User:      "postgres",
		Password:  "postgres",
		Port:      "5451",
		DbName:    "postgres",
		SSLMode:   "disable",
		TimeZone:  "America/Sao_Paulo",
	})
	assert.NoError(t, err)

	var actualUsername, actualEmail, userID string
	var deleted *time.Time
	err = connect.Raw("SELECT username, email, deleted_at, id FROM users WHERE username = ? AND email = ?", username, email).
		Row().
		Scan(&actualUsername, &actualEmail, &deleted, &userID)
	assert.NoError(t, err)

	assert.Equal(t, username, actualUsername)
	assert.Equal(t, email, actualEmail)
	assert.NotNil(t, deleted)

	err = connect.Exec("DELETE FROM users WHERE id = ?", userID).Error
	assert.NoError(t, err)
}
