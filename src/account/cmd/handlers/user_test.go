package handlers

import (
	"account/internal/user"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Create(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockRepository) Update(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockRepository) FindById(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_CreateUser(t *testing.T) {
	cases := []struct {
		name     string
		req      string
		code     int
		mockFunc func(repository *mockRepository)
		want     string
	}{
		{
			name: "success",
			code: http.StatusCreated,
			req:  `{"username":"example_user","email":"user@example.com"}`,
			mockFunc: func(repository *mockRepository) {
				repository.On("Create", mock.Anything, &user.User{
					Username: "example_user",
					Email:    "user@example.com",
				}).Return(nil).Once()
			},
			want: `{"created_at":"0001-01-01T00:00:00Z","email":"user@example.com","id":"","updated_at":"0001-01-01T00:00:00Z","username":"example_user"}`,
		},
		{
			name: "failed because error in delete",
			code: http.StatusInternalServerError,
			req:  `{"username":"example_user","email":"user@example.com"}`,
			mockFunc: func(repository *mockRepository) {
				repository.On("Create", mock.Anything, &user.User{
					Username: "example_user",
					Email:    "user@example.com",
				}).Return(errors.New("mock-error")).Once()
			},
			want: `{"message":"mock-error"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mockRepository)

			tc.mockFunc(repo)

			s := NewUserHandler(repo)

			app := fiber.New()
			app.Post("/users", s.CreateUser)

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(tc.req))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, 1)
			assert.NoError(t, err)

			defer resp.Body.Close()

			var want map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&want)
			assert.NoError(t, err)

			marshal, err := json.Marshal(want)
			assert.NoError(t, err)

			assert.Equal(t, tc.want, string(marshal))
			assert.Equal(t, tc.code, resp.StatusCode)
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	cases := []struct {
		name     string
		req      string
		code     int
		mockFunc func(repository *mockRepository)
		want     string
	}{
		{
			name: "success",
			code: http.StatusCreated,
			req:  `{"username":"example_user","email":"user@example.com"}`,
			mockFunc: func(repository *mockRepository) {
				repository.On("Update", mock.Anything, &user.User{
					ID:       "1",
					Username: "example_user",
					Email:    "user@example.com",
				}).Return(nil).Once()
			},
			want: `{"created_at":"0001-01-01T00:00:00Z","email":"user@example.com","id":"1","updated_at":"0001-01-01T00:00:00Z","username":"example_user"}`,
		},
		{
			name: "failed because error in update",
			code: http.StatusInternalServerError,
			req:  `{"username":"example_user","email":"user@example.com"}`,
			mockFunc: func(repository *mockRepository) {
				repository.On("Update", mock.Anything, &user.User{
					ID:       "1",
					Username: "example_user",
					Email:    "user@example.com",
				}).Return(errors.New("mock-error")).Once()
			},
			want: `{"message":"mock-error"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mockRepository)

			tc.mockFunc(repo)

			s := NewUserHandler(repo)

			app := fiber.New()
			app.Put("/users/:id", s.UpdateUser)

			req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBufferString(tc.req))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, 1)
			assert.NoError(t, err)

			defer resp.Body.Close()

			var want map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&want)
			assert.NoError(t, err)

			marshal, err := json.Marshal(want)
			assert.NoError(t, err)

			assert.Equal(t, tc.want, string(marshal))
			assert.Equal(t, tc.code, resp.StatusCode)
		})
	}
}

func TestService_FindUser(t *testing.T) {
	cases := []struct {
		name     string
		code     int
		mockFunc func(repository *mockRepository)
		want     string
	}{
		{
			name: "success",
			code: http.StatusCreated,
			mockFunc: func(repository *mockRepository) {
				repository.On("FindById", mock.Anything, "1").Return(&user.User{
					Username: "example_user",
					Email:    "user@example.com",
				}, nil).Once()
			},
			want: `{"created_at":"0001-01-01T00:00:00Z","email":"user@example.com","id":"","updated_at":"0001-01-01T00:00:00Z","username":"example_user"}`,
		},
		{
			name: "failed because error in find by id",
			code: http.StatusInternalServerError,
			mockFunc: func(repository *mockRepository) {
				repository.On("FindById", mock.Anything, "1").
					Return((*user.User)(nil), errors.New("mock-error")).Once()
			},
			want: `{"message":"mock-error"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mockRepository)

			tc.mockFunc(repo)

			s := NewUserHandler(repo)

			app := fiber.New()
			app.Get("/users/:id", s.FindUser)

			req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, 1)
			assert.NoError(t, err)

			defer resp.Body.Close()

			var want map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&want)
			assert.NoError(t, err)

			marshal, err := json.Marshal(want)
			assert.NoError(t, err)

			assert.Equal(t, tc.want, string(marshal))
			assert.Equal(t, tc.code, resp.StatusCode)
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	cases := []struct {
		name     string
		code     int
		mockFunc func(repository *mockRepository)
		want     string
	}{
		{
			name: "success",
			code: http.StatusCreated,
			mockFunc: func(repository *mockRepository) {
				repository.On("Delete", mock.Anything, "1").Return(nil).Once()
			},
			want: `{"message":"success"}`,
		},
		{
			name: "failed because error in find by id",
			code: http.StatusInternalServerError,
			mockFunc: func(repository *mockRepository) {
				repository.On("Delete", mock.Anything, "1").
					Return(errors.New("mock-error")).Once()
			},
			want: `{"message":"mock-error"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mockRepository)

			tc.mockFunc(repo)

			s := NewUserHandler(repo)

			app := fiber.New()
			app.Delete("/users/:id", s.DeleteUser)

			req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, 1)
			assert.NoError(t, err)

			defer resp.Body.Close()

			var want map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&want)
			assert.NoError(t, err)

			marshal, err := json.Marshal(want)
			assert.NoError(t, err)

			assert.Equal(t, tc.want, string(marshal))
			assert.Equal(t, tc.code, resp.StatusCode)
		})
	}
}
