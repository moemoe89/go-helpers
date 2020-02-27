package helpers_test

import (
	"github.com/moemoe89/go-helpers"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderByHandler(t *testing.T) {
	selectField := "id"
	orderBy := helpers.OrderByHandler(selectField, "db", UserModel{})
	expectedOrderBy := "id ASC"

	assert.Equal(t, expectedOrderBy, orderBy)
}

func TestOrderByHandlerDesc(t *testing.T) {
	selectField := "-id"
	orderBy := helpers.OrderByHandler(selectField, "db", UserModel{})
	expectedOrderBy := "id DESC"

	assert.Equal(t, expectedOrderBy, orderBy)
}

func TestOrderByHandleEmptyField(t *testing.T) {
	selectField := ""
	orderBy := helpers.OrderByHandler(selectField, "db", UserModel{})
	expectedOrderBy := ""

	assert.Equal(t, expectedOrderBy, orderBy)
}

func TestOrderByHandleNotFoundField(t *testing.T) {
	selectField := "x"
	orderBy := helpers.OrderByHandler(selectField, "db", UserModel{})
	expectedOrderBy := ""

	assert.Equal(t, expectedOrderBy, orderBy)
}
