package helpers

import (
	"errors"
	"strconv"
)

// PaginationSetter represent the helpers for pagination
func PaginationSetter(perPageString, pageString string) (int, int, int, error) {
	if len(perPageString) == 0 {
		perPageString = "10"
	}
	perPage, err := strconv.Atoi(perPageString)
	if err != nil {
		return 0, 0, 0, errors.New("Invalid parameter per_page: not an int")
	}

	if len(pageString) == 0 {
		pageString = "1"
	}
	page, err := strconv.Atoi(pageString)
	if err != nil {
		return 0, 0, 0, errors.New("Invalid parameter page: not an int")
	}

	showPage := page
	if showPage < 1 {
		showPage = 1
		page = 1
	}
	offset := (page - 1) * perPage

	return offset, perPage, showPage, nil
}
