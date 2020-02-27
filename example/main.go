//
//  main.go
//  helpers
//
//  Copyright © 2020. All rights reserved.
//

package main

import (
	"github.com/moemoe89/go-helpers"

	"fmt"
	"time"
)

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

func main() {
	// check in tag
	filterField := "id,name,email,phone,address,created_at,updated_at"
	res := helpers.CheckInTag(UserModel{}, filterField, "db")
	fmt.Println(res)

	// check match tag
	selectField := "id"
	field := helpers.CheckMatchTag(UserModel{}, selectField, "db")
	fmt.Println(field)

	// handling order by
	orderBy := helpers.OrderByHandler(selectField, "db", UserModel{})
	fmt.Println(orderBy)

	// get pagination setter
	offset, perPage, showPage, err := helpers.PaginationSetter("10", "0")
	if err != nil {
		panic(err)
	}
	fmt.Println(offset)
	fmt.Println(perPage)
	fmt.Println(showPage)
}