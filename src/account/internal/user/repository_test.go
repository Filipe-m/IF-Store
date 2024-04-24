package user

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestRepository_Create(t *testing.T) {
	cases := []struct {
		name     string
		req      *User
		mockFunc func(sqlMock sqlmock.Sqlmock)
		wantErr  error
	}{
		{
			name: "should return success with error nil",
			req: &User{
				Username:  "example_user",
				Email:     "user@example.com",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: gorm.DeletedAt{},
			},
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				sqlMock.ExpectQuery("INSERT INTO").
					WithArgs(
						"example_user",
						"user@example.com",
						time.Time{},
						time.Time{},
						nil,
					).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
				sqlMock.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name: "should return failed with error",
			req: &User{
				Username:  "example_user",
				Email:     "user@example.com",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: gorm.DeletedAt{},
			},
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				sqlMock.ExpectQuery("INSERT INTO").
					WithArgs(
						"example_user",
						"user@example.com",
						time.Time{},
						time.Time{},
						nil,
					).WillReturnError(errors.New("error"))
				sqlMock.ExpectRollback()
			},
			wantErr: errors.New("error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			conn, mocksql, err := sqlmock.New()
			assert.Nil(t, err)

			defer conn.Close()

			dialector := postgres.New(postgres.Config{
				DSN:                  "sqlmock_db_0",
				DriverName:           "postgres",
				Conn:                 conn,
				PreferSimpleProtocol: true,
			})

			db, err := gorm.Open(dialector, &gorm.Config{
				NowFunc: func() time.Time {
					return time.Time{}
				},
			})
			assert.NoError(t, err)

			repos := NewRepository(db)
			tc.mockFunc(mocksql)

			ctx := context.Background()

			err = repos.Create(ctx, tc.req)
			assert.Equal(t, tc.wantErr, err)
			assert.Nil(t, mocksql.ExpectationsWereMet())
		})
	}
}

func TestRepository_Update(t *testing.T) {
	cases := []struct {
		name     string
		req      *User
		mockFunc func(sqlMock sqlmock.Sqlmock)
		wantErr  error
	}{
		{
			name: "should return success with error nil",
			req: &User{
				ID:        "1",
				Username:  "example_user",
				Email:     "user@example.com",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: gorm.DeletedAt{},
			},
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec("UPDATE").
					WithArgs(
						"example_user",
						"user@example.com",
						"1",
					).WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit()

				sqlMock.ExpectQuery(`SELECT`).WillReturnRows(
					sqlmock.NewRows([]string{
						"id",
						"username",
						"email",
						"created_at",
						"updated_at",
						"deleted_at",
					}).AddRow(
						"1",
						"example_user",
						"user@example.com",
						time.Time{},
						time.Time{},
						nil,
					),
				)
			},
			wantErr: nil,
		},
		{
			name: "should failed because error in select",
			req: &User{
				ID:        "1",
				Username:  "example_user",
				Email:     "user@example.com",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: gorm.DeletedAt{},
			},
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec("UPDATE").
					WithArgs(
						"example_user",
						"user@example.com",
						"1",
					).WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit()

				sqlMock.ExpectQuery(`SELECT`).WillReturnError(errors.New("error"))
			},
			wantErr: errors.New("error"),
		},
		{
			name: "should failed because error in update",
			req: &User{
				ID:        "1",
				Username:  "example_user",
				Email:     "user@example.com",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: gorm.DeletedAt{},
			},
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec("UPDATE").
					WithArgs(
						"example_user",
						"user@example.com",
						"1",
					).WillReturnError(errors.New("error"))
				sqlMock.ExpectRollback()
			},
			wantErr: errors.New("error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			conn, mocksql, err := sqlmock.New()
			assert.Nil(t, err)

			defer conn.Close()

			dialector := postgres.New(postgres.Config{
				DSN:                  "sqlmock_db_0",
				DriverName:           "postgres",
				Conn:                 conn,
				PreferSimpleProtocol: true,
			})

			db, err := gorm.Open(dialector, &gorm.Config{
				NowFunc: func() time.Time {
					return time.Time{}
				},
			})
			assert.NoError(t, err)

			repos := NewRepository(db)
			tc.mockFunc(mocksql)

			ctx := context.Background()

			err = repos.Update(ctx, tc.req)
			assert.Equal(t, tc.wantErr, err)
			assert.Nil(t, mocksql.ExpectationsWereMet())
		})
	}
}

func TestRepository_FindById(t *testing.T) {
	cases := []struct {
		name     string
		req      string
		mockFunc func(sqlMock sqlmock.Sqlmock)
		want     *User
		wantErr  error
	}{
		{
			name: "should return success with error nil",
			req:  "1",
			want: &User{
				ID:        "1",
				Username:  "example_user",
				Email:     "user@example.com",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: gorm.DeletedAt{
					Time:  time.Time{},
					Valid: false,
				},
			},
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery(`SELECT`).WillReturnRows(
					sqlmock.NewRows([]string{
						"id",
						"username",
						"email",
						"created_at",
						"updated_at",
						"deleted_at",
					}).AddRow(
						"1",
						"example_user",
						"user@example.com",
						time.Time{},
						time.Time{},
						nil,
					),
				)
			},
			wantErr: nil,
		},
		{
			name: "should return failed because error in select",
			req:  "1",
			want: nil,
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery(`SELECT`).WillReturnError(errors.New("error"))
			},
			wantErr: errors.New("error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			conn, mocksql, err := sqlmock.New()
			assert.Nil(t, err)

			defer conn.Close()

			dialector := postgres.New(postgres.Config{
				DSN:                  "sqlmock_db_0",
				DriverName:           "postgres",
				Conn:                 conn,
				PreferSimpleProtocol: true,
			})

			db, err := gorm.Open(dialector, &gorm.Config{
				NowFunc: func() time.Time {
					return time.Time{}
				},
			})
			assert.NoError(t, err)

			repos := NewRepository(db)
			tc.mockFunc(mocksql)

			ctx := context.Background()

			want, err := repos.FindById(ctx, tc.req)
			assert.Equal(t, tc.want, want)
			assert.Equal(t, tc.wantErr, err)
			assert.Nil(t, mocksql.ExpectationsWereMet())
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	cases := []struct {
		name     string
		req      string
		mockFunc func(sqlMock sqlmock.Sqlmock)
		wantErr  error
	}{
		{
			name: "should return success with error nil",
			req:  "1",
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec("UPDATE").
					WithArgs(
						time.Time{},
						"1",
					).WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name: "should return success with error nil",
			req:  "1",
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec("UPDATE").
					WithArgs(
						time.Time{},
						"1",
					).WillReturnError(errors.New("error"))
				sqlMock.ExpectRollback()
			},
			wantErr: errors.New("error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			conn, mocksql, err := sqlmock.New()
			assert.Nil(t, err)

			defer conn.Close()

			dialector := postgres.New(postgres.Config{
				DSN:                  "sqlmock_db_0",
				DriverName:           "postgres",
				Conn:                 conn,
				PreferSimpleProtocol: true,
			})

			db, err := gorm.Open(dialector, &gorm.Config{
				NowFunc: func() time.Time {
					return time.Time{}
				},
			})
			assert.NoError(t, err)

			repos := NewRepository(db)
			tc.mockFunc(mocksql)

			ctx := context.Background()

			err = repos.Delete(ctx, tc.req)
			assert.Equal(t, tc.wantErr, err)
			assert.Nil(t, mocksql.ExpectationsWereMet())
		})
	}
}
