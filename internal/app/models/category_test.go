package models

import (
	"testing"

	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetCategory(t *testing.T) {
	mockDB, mock, sqlxDB := MockDB(t)
	defer mockDB.Close()

	var cols []string = []string{"id", "alias", "name"}
	mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).
		AddRow(1, "other", "Other"))

	cm := NewCategoryModel(sqlxDB)
	c, ok := cm.Get(1)

	expect := &Category{
		ID:    1,
		Alias: "other",
		Name:  "Other",
	}

	require.True(t, ok)
	require.Equal(t, expect, c)
}

func TestGetCategories(t *testing.T) {
	mockDB, mock, sqlxDB := MockDB(t)
	defer mockDB.Close()

	var cols []string = []string{"id", "alias", "name"}
	mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).
		AddRow(1, "other", "Other").AddRow(2, "some_other", "Some Other"))

	cm := NewCategoryModel(sqlxDB)
	c, err := cm.GetAll()
	if err != nil {
		t.Fail()
	}

	expect := []Category{
		{
			ID:    1,
			Alias: "other",
			Name:  "Other",
		},
		{
			ID:    2,
			Alias: "some_other",
			Name:  "Some Other",
		},
	}
	require.Equal(t, expect, c)
}
