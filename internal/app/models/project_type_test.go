package models

import (
	"testing"

	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetProjectTypes(t *testing.T) {
	mockDB, mock, sqlxDB := MockDB(t)
	defer mockDB.Close()

	var cols []string = []string{"id", "alias", "name"}
	mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).
		AddRow(1, "other", "Other").
		AddRow(2, "some_other", "Some Other"))

	ptm := NewProjectTypeModel(sqlxDB)
	pt, err := ptm.GetAll()
	if err != nil {
		t.Fail()
	}

	expect := []ProjectType{
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
	require.Equal(t, expect, pt)
}
