package helpers_test

import (
	"github.com/moemoe89/go-helpers"

	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const UserSelectField = "id,name,email,phone,address,created_at,updated_at"

type UserModel struct {
	ID        string     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Email     string     `json:"email" db:"email"`
	Phone     string     `json:"phone" db:"phone"`
	Address   string     `json:"address" db:"address"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

func TestCheckInDBAtt(t *testing.T) {
	selectField := ""
	filterField := UserSelectField
	res := helpers.CheckInTag(UserModel{}, filterField, "db")
	if len(res) > 0 {
		selectField = strings.Join(res, ",")
	}

	assert.Equal(t, selectField, filterField)
}

func TestCheckMatchDBAtt(t *testing.T) {
	selectField := "id"
	field := helpers.CheckMatchTag(UserModel{}, selectField, "db")

	assert.Equal(t, selectField, field)
}
