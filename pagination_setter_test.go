package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginationSetterFailPerPage(t *testing.T) {
	_, _, _, err := PaginationSetter("a", "")
	expectedMsg := "Invalid parameter per_page: not an int"

	assert.Equal(t, expectedMsg, err.Error())
}

func TestPaginationSetterFailPage(t *testing.T) {
	_, _, _, err := PaginationSetter("", "b")
	expectedMsg := "Invalid parameter page: not an int"

	assert.Equal(t, expectedMsg, err.Error())
}

func TestPaginationSetterPage0(t *testing.T) {
	_, _, _, err := PaginationSetter("10", "0")

	assert.Equal(t, nil, err)
}

func TestPaginationSetter(t *testing.T) {
	_, _, _, err := PaginationSetter("", "")

	assert.Equal(t, nil, err)
}
