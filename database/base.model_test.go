package database

import (
	"github.com/usercoredev/usercore/utils"
	"testing"
)

func TestPageMetadata_ConvertToOrder_UnsafeValues(t *testing.T) {
	user := User{}
	pm := &utils.PageMetadata{OrderBy: "1; DROP TABLE users", Order: ""}
	orderClause := user.ConvertToOrder(*pm)
	if orderClause != "created_at desc" { // Assuming you sanitize or ignore unsafe values
		t.Errorf("Expected default order clause for unsafe OrderBy, got %s", orderClause)
	}
}
