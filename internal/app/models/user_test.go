package models

import (
	"testing"

	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestFindByID(t *testing.T) {
	mockDB, mock, sqlxDB := MockDB(t)
	defer mockDB.Close()

	var cols []string = []string{"id", "username", "first_name", "last_name", "avatar", "email", "is_admin", "project_count", "success_rate"}
	mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).
		AddRow(1, "johnny86", "John", "Doe", "https://avatar.com", "john@gmail.com", false, 1, 100.0))

	um := NewUserModel(sqlxDB)
	u, ok := um.FindByID(1)

	expect := User{
		ID:   1,
		Username: "johnny86",
		FirstName: "John",
		LastName: "Doe",
		Avatar: "https://avatar.com",
		Email: "john@gmail.com",
		ProjectCount: 1,
		SuccessRate: 100.0,
	}

	require.True(t, ok)
	require.Equal(t, expect, u)
}

func TestCreate(t *testing.T) {
	mockDB, mock, sqlxDB := MockDB(t)
	defer mockDB.Close()

	expect := User{
		ID:   1,
		Username: "johnny86",
		FirstName: "John",
		LastName: "Doe",
		Avatar: "https://avatar.com",
		Email: "john@gmail.com",
	}
	mock.ExpectExec("INSERT INTO users").WithArgs(
		1, "johnny86", "John", "Doe", "https://avatar.com", "john@gmail.com",
		).WillReturnResult(sqlmock.NewResult(1, 1))
	um := NewUserModel(sqlxDB)
	_, err := um.Create(&expect)

	require.NoError(t, err)
}
