package helpers_test

import (
	"github.com/moemoe89/go-helpers"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginationSetterFailPerPage(t *testing.T) {
	_, _, _, err := helpers.PaginationSetter("a", "")
	expectedMsg := "Invalid parameter per_page: not an int"

	assert.Equal(t, expectedMsg, err.Error())
}

func TestPaginationSetterFailPage(t *testing.T) {
	_, _, _, err := helpers.PaginationSetter("", "b")
	expectedMsg := "Invalid parameter page: not an int"

	assert.Equal(t, expectedMsg, err.Error())
}

func TestPaginationSetterPage0(t *testing.T) {
	_, _, _, err := helpers.PaginationSetter("10", "0")

	assert.Equal(t, nil, err)
}

func TestPaginationSetter(t *testing.T) {
	_, _, _, err := helpers.PaginationSetter("", "")

	assert.Equal(t, nil, err)
}
