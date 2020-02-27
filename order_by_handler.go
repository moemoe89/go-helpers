package helpers

import (
	"fmt"
	"strings"
)

// OrderByHandler represent the helpers for handle order by query
func OrderByHandler(field, tag string, model interface{}) string {
	if field == "" {
		return field
	}

	sort := "ASC"
	splitDesc := strings.Split(field, "")
	if len(splitDesc) > 0 {
		if splitDesc[0] == "-" {
			sort = "DESC"
		}
	}

	field = strings.Replace(field, "-", "", -1)
	field = CheckMatchTag(model, field, tag)

	if field == "" {
		return field
	}

	return fmt.Sprintf("%s %s", field, sort)
}
